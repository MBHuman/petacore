#!/usr/bin/env bash
set -euo pipefail

PORT="${PORT:-5432}"
STORE="${STORE:-inmemory}"

LOG_FILE="/tmp/pcore_server_${PORT}.log"
WRAPPER_PID=""

is_listening() {
  (echo >"/dev/tcp/127.0.0.1/${PORT}") >/dev/null 2>&1
}

port_pid() {
  # Возвращает PID процесса, который LISTEN на порту (если есть)
  # -t: только PID, -sTCP:LISTEN: только LISTEN
  lsof -nP -t -iTCP:"${PORT}" -sTCP:LISTEN 2>/dev/null || true
}

cleanup() {
  echo "🧹 Stopping server..."

  # 1) Сначала прибиваем того, кто реально слушает порт
  local pids
  pids="$(port_pid)"
  if [[ -n "${pids}" ]]; then
    echo "🔪 Killing listener on port ${PORT}: ${pids}"
    kill -TERM ${pids} 2>/dev/null || true
    sleep 0.2
    kill -KILL ${pids} 2>/dev/null || true
  fi

  # 2) На всякий случай прибиваем wrapper go run
  if [[ -n "${WRAPPER_PID}" ]] && kill -0 "${WRAPPER_PID}" 2>/dev/null; then
    kill -TERM "${WRAPPER_PID}" 2>/dev/null || true
  fi

  # 3) Ждём освобождения порта
  local deadline=$((SECONDS + 60))
  while is_listening && (( SECONDS < deadline )); do
    sleep 0.1
  done

  if is_listening; then
    echo "⚠️  Port ${PORT} is still in use. Holder:"
    lsof -nP -iTCP:"${PORT}" -sTCP:LISTEN || true
  fi
}

trap cleanup EXIT INT TERM

# Pre-check
if is_listening; then
  echo "❌ Port ${PORT} is already in use. Holder:"
  lsof -nP -iTCP:"${PORT}" -sTCP:LISTEN || true
  echo "   Try: PORT=55432 ./scripts/run-tests.sh"
  exit 1
fi

echo "🚀 Starting pcore server (store=${STORE}, port=${PORT})..."
echo "📝 Logging to ${LOG_FILE}"

go run cmd/pcore/main.go -store "${STORE}" >"${LOG_FILE}" 2>&1 &
WRAPPER_PID=$!

echo "⏳ Waiting for server to listen on ${PORT}..."
deadline=$((SECONDS + 60))
while (( SECONDS < deadline )); do
  if ! kill -0 "${WRAPPER_PID}" 2>/dev/null; then
    echo "❌ Server wrapper exited during startup."
    echo "---- server log ----"
    tail -n 200 "${LOG_FILE}" || true
    exit 1
  fi

  if is_listening; then
    echo "✅ Server is listening"
    break
  fi

  sleep 0.2
done

if ! is_listening; then
  echo "❌ Server did not open port ${PORT} in time."
  echo "---- server log ----"
  tail -n 200 "${LOG_FILE}" || true
  exit 1
fi

echo "🧪 Running SQL tests..."
go test -count=1 ./tests/sql_test

echo "✅ Tests finished"
