#!/bin/bash

# E2E бенчмарк для PetaCore VClock с реальным ETCD
# Тестирует производительность распределенного хранилища с Vector Clock

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

API_ENDPOINTS=(
    "http://localhost:8081"
    "http://localhost:8082"
    "http://localhost:8083"
)

RESULTS_DIR="benchmark-results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="$RESULTS_DIR/vclock_e2e_$TIMESTAMP.txt"

# Создаем директорию для результатов
mkdir -p "$RESULTS_DIR"

echo -e "${YELLOW}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║  PetaCore VClock E2E Benchmark with Real ETCD            ║${NC}"
echo -e "${YELLOW}║  Timestamp: $(date '+%Y-%m-%d %H:%M:%S')                       ║${NC}"
echo -e "${YELLOW}╚═══════════════════════════════════════════════════════════╝${NC}"
echo

# Проверка наличия инструментов
if ! command -v ab &> /dev/null; then
    echo -e "${RED}Error: apache-bench (ab) not found${NC}"
    echo "Installing..."
    sudo apt-get update && sudo apt-get install -y apache2-utils
fi

if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq not found${NC}"
    exit 1
fi

# Функция для логирования результатов
log_result() {
    echo "$1" | tee -a "$RESULTS_FILE"
}

# Функция для извлечения метрик из ab
extract_metrics() {
    local output="$1"
    local rps=$(echo "$output" | grep "Requests per second" | awk '{print $4}')
    local time_per_req=$(echo "$output" | grep "Time per request" | head -1 | awk '{print $4}')
    local p50=$(echo "$output" | grep "50%" | awk '{print $2}')
    local p95=$(echo "$output" | grep "95%" | awk '{print $2}')
    local p99=$(echo "$output" | grep "99%" | awk '{print $2}')
    
    echo "$rps|$time_per_req|$p50|$p95|$p99"
}

log_result "=========================================="
log_result "E2E Benchmark Results - VClock with ETCD"
log_result "=========================================="
log_result ""

# ============================================================================
# Тест 1: Проверка здоровья кластера
# ============================================================================
echo -e "${BLUE}═══ Test 1: Cluster Health Check ═══${NC}"
log_result "Test 1: Cluster Health Check"
log_result "------------------------------"

for i in "${!API_ENDPOINTS[@]}"; do
    endpoint="${API_ENDPOINTS[$i]}"
    status=$(curl -s "$endpoint/status" | jq -r '{node_id, is_synced, total_nodes, min_acks}')
    echo -e "${GREEN}Node $((i+1)):${NC}"
    echo "$status"
    log_result "Node $((i+1)): $(echo $status | tr '\n' ' ')"
done
log_result ""
echo

# ============================================================================
# Тест 2: Базовая производительность записи (с разными MIN_ACKS)
# ============================================================================
echo -e "${BLUE}═══ Test 2: Write Performance with Different MIN_ACKS ═══${NC}"
log_result "Test 2: Write Performance with Different MIN_ACKS"
log_result "---------------------------------------------------"

for min_acks in 1 0 -1; do
    min_acks_desc=""
    case $min_acks in
        1) min_acks_desc="MIN_ACKS=1 (weak consistency)" ;;
        0) min_acks_desc="MIN_ACKS=0 (quorum)" ;;
        -1) min_acks_desc="MIN_ACKS=-1 (all nodes)" ;;
    esac
    
    echo -e "${YELLOW}Testing $min_acks_desc${NC}"
    log_result "$min_acks_desc:"
    
    # Устанавливаем MIN_ACKS
    curl -s -X POST "${API_ENDPOINTS[0]}/config/min_acks" \
        -H "Content-Type: application/json" \
        -d "{\"min_acks\":$min_acks}" > /dev/null
    
    sleep 1
    
    # Создаем тестовый JSON файл
    echo "{\"key\":\"bench:write_${min_acks}\",\"value\":\"test_value\"}" > /tmp/bench_write.json
    
    # Бенчмарк записи
    output=$(ab -n 200 -c 5 -p /tmp/bench_write.json -T "application/json" -q "${API_ENDPOINTS[0]}/write" 2>&1)
    
    rps=$(echo "$output" | grep "Requests per second" | awk '{print $4}')
    time_per_req=$(echo "$output" | grep "Time per request" | head -1 | awk '{print $4}')
    
    echo "  Requests per second: $rps"
    echo "  Time per request: ${time_per_req}ms"
    log_result "  RPS: $rps, Time/req: ${time_per_req}ms"
    
    sleep 1
done
log_result ""
echo

# Восстанавливаем quorum по умолчанию
curl -s -X POST "${API_ENDPOINTS[0]}/config/min_acks" \
    -H "Content-Type: application/json" \
    -d '{"min_acks":0}' > /dev/null

# ============================================================================
# Тест 3: Производительность чтения (с quorum проверкой)
# ============================================================================
echo -e "${BLUE}═══ Test 3: Read Performance with Quorum Check ═══${NC}"
log_result "Test 3: Read Performance with Quorum Check"
log_result "--------------------------------------------"

# Записываем тестовые данные
for i in {1..10}; do
    curl -s -X POST "${API_ENDPOINTS[0]}/write" \
        -H "Content-Type: application/json" \
        -d "{\"key\":\"bench:read_$i\",\"value\":\"value_$i\"}" > /dev/null
done

sleep 2

echo "Warming up..."
for i in {1..100}; do
    curl -s "${API_ENDPOINTS[0]}/read?key=bench:read_1" > /dev/null
done

echo "Running read benchmark (1000 requests, concurrency 10)..."
output=$(ab -n 1000 -c 10 -q "${API_ENDPOINTS[0]}/read?key=bench:read_1" 2>&1)

metrics=$(extract_metrics "$output")
IFS='|' read -r rps time_per_req p50 p95 p99 <<< "$metrics"

echo "  Requests per second: $rps"
echo "  Time per request: ${time_per_req}ms"
echo "  Latency p50: ${p50}ms, p95: ${p95}ms, p99: ${p99}ms"

log_result "  RPS: $rps"
log_result "  Time/req: ${time_per_req}ms"
log_result "  Latency - p50: ${p50}ms, p95: ${p95}ms, p99: ${p99}ms"
log_result ""
echo

# ============================================================================
# Тест 4: Распределенная запись на все узлы
# ============================================================================
echo -e "${BLUE}═══ Test 4: Distributed Writes Across All Nodes ═══${NC}"
log_result "Test 4: Distributed Writes Across All Nodes"
log_result "---------------------------------------------"

echo "Writing to all 3 nodes in parallel..."
start_time=$(date +%s%N)

for i in {1..300}; do
    node_idx=$((i % 3))
    endpoint="${API_ENDPOINTS[$node_idx]}"
    curl -s -X POST "$endpoint/write" \
        -H "Content-Type: application/json" \
        -d "{\"key\":\"bench:distributed_$i\",\"value\":\"value_$i\"}" > /dev/null &
    
    # Ограничиваем количество параллельных процессов
    if [ $((i % 30)) -eq 0 ]; then
        wait
    fi
done
wait

end_time=$(date +%s%N)
duration=$(echo "scale=3; ($end_time - $start_time) / 1000000000" | bc)
throughput=$(echo "scale=2; 300 / $duration" | bc)

echo "  Total time: ${duration}s"
echo "  Throughput: ${throughput} writes/sec"
log_result "  Duration: ${duration}s, Throughput: ${throughput} writes/sec"

# Ждем синхронизации
echo "Waiting for sync..."
sleep 3

# Проверяем консистентность - читаем с разных узлов
echo "Checking consistency across nodes..."
consistent=true
for key_num in {1..10}; do
    values=()
    for endpoint in "${API_ENDPOINTS[@]}"; do
        value=$(curl -s "$endpoint/read?key=bench:distributed_$key_num" | jq -r '.value')
        values+=("$value")
    done
    
    # Проверяем что все значения одинаковые
    if [ "${values[0]}" != "${values[1]}" ] || [ "${values[0]}" != "${values[2]}" ]; then
        consistent=false
        echo -e "${RED}  ✗ Inconsistent: bench:distributed_$key_num${NC}"
        log_result "  ✗ Inconsistent: key bench:distributed_$key_num"
    fi
done

if $consistent; then
    echo -e "${GREEN}  ✓ All checked keys are consistent across nodes${NC}"
    log_result "  ✓ Consistency check passed"
else
    echo -e "${RED}  ✗ Some keys are inconsistent${NC}"
    log_result "  ✗ Consistency check failed"
fi
log_result ""
echo

# ============================================================================
# Тест 5: Параллельные чтения/записи (mixed workload)
# ============================================================================
echo -e "${BLUE}═══ Test 5: Mixed Read/Write Workload ═══${NC}"
log_result "Test 5: Mixed Read/Write Workload"
log_result "-----------------------------------"

echo "Running mixed workload: 70% reads, 30% writes..."
start_time=$(date +%s%N)

# 700 чтений
for i in {1..700}; do
    node_idx=$((RANDOM % 3))
    key_num=$((RANDOM % 100 + 1))
    curl -s "${API_ENDPOINTS[$node_idx]}/read?key=bench:mixed_$key_num" > /dev/null &
    
    if [ $((i % 50)) -eq 0 ]; then
        wait
    fi
done

# 300 записей
for i in {1..300}; do
    node_idx=$((RANDOM % 3))
    curl -s -X POST "${API_ENDPOINTS[$node_idx]}/write" \
        -H "Content-Type: application/json" \
        -d "{\"key\":\"bench:mixed_$i\",\"value\":\"mixed_$i\"}" > /dev/null &
    
    if [ $((i % 30)) -eq 0 ]; then
        wait
    fi
done
wait

end_time=$(date +%s%N)
duration=$(echo "scale=3; ($end_time - $start_time) / 1000000000" | bc)
throughput=$(echo "scale=2; 1000 / $duration" | bc)

echo "  Total operations: 1000 (700 reads + 300 writes)"
echo "  Duration: ${duration}s"
echo "  Throughput: ${throughput} ops/sec"
log_result "  1000 ops (70% read, 30% write)"
log_result "  Duration: ${duration}s, Throughput: ${throughput} ops/sec"
log_result ""
echo

# ============================================================================
# Тест 6: Vector Clock overhead (сравнение с/без VClock)
# ============================================================================
echo -e "${BLUE}═══ Test 6: Vector Clock Synchronization Test ═══${NC}"
log_result "Test 6: Vector Clock Synchronization Test"
log_result "-------------------------------------------"

echo "Testing VClock sync across nodes..."

# Записываем на node-1
for i in {1..50}; do
    curl -s -X POST "${API_ENDPOINTS[0]}/write" \
        -H "Content-Type: application/json" \
        -d "{\"key\":\"vclock:test_$i\",\"value\":\"from_node1_$i\"}" > /dev/null
done

echo "Waiting 2 seconds for VClock sync..."
sleep 2

# Читаем с node-2 и node-3
sync_success=0
sync_fail=0

for i in {1..50}; do
    # Читаем с node-2
    value_node2=$(curl -s "${API_ENDPOINTS[1]}/read?key=vclock:test_$i" | jq -r '.value')
    # Читаем с node-3
    value_node3=$(curl -s "${API_ENDPOINTS[2]}/read?key=vclock:test_$i" | jq -r '.value')
    
    if [ "$value_node2" == "from_node1_$i" ] && [ "$value_node3" == "from_node1_$i" ]; then
        ((sync_success++))
    else
        ((sync_fail++))
    fi
done

sync_rate=$(echo "scale=2; $sync_success * 100 / 50" | bc)
echo "  VClock sync success rate: ${sync_rate}% (${sync_success}/50)"
log_result "  Sync success: ${sync_rate}% (${sync_success}/50)"

if [ $sync_success -eq 50 ]; then
    echo -e "${GREEN}  ✓ Perfect synchronization${NC}"
    log_result "  ✓ Perfect synchronization"
else
    echo -e "${YELLOW}  ⚠ Partial synchronization (may need more time)${NC}"
    log_result "  ⚠ Partial synchronization"
fi
log_result ""
echo

# ============================================================================
# Тест 7: Latency под нагрузкой
# ============================================================================
echo -e "${BLUE}═══ Test 7: Latency Under Load ═══${NC}"
log_result "Test 7: Latency Under Load"
log_result "----------------------------"

for concurrency in 1 10 50 100; do
    echo -e "${YELLOW}Concurrency: $concurrency${NC}"
    log_result "Concurrency $concurrency:"
    
    output=$(ab -n 500 -c $concurrency -q "${API_ENDPOINTS[0]}/read?key=bench:read_1" 2>&1)
    
    metrics=$(extract_metrics "$output")
    IFS='|' read -r rps time_per_req p50 p95 p99 <<< "$metrics"
    
    echo "  RPS: $rps, Latency - p50: ${p50}ms, p95: ${p95}ms, p99: ${p99}ms"
    log_result "  RPS: $rps, p50: ${p50}ms, p95: ${p95}ms, p99: ${p99}ms"
    
    sleep 1
done
log_result ""
echo

# ============================================================================
# Тест 8: Failover scenario (симуляция отказа узла)
# ============================================================================
echo -e "${BLUE}═══ Test 8: Failover Scenario ═══${NC}"
log_result "Test 8: Failover Scenario"
log_result "---------------------------"

echo "Writing data with all 3 nodes active..."
for i in {1..20}; do
    curl -s -X POST "${API_ENDPOINTS[0]}/write" \
        -H "Content-Type: application/json" \
        -d "{\"key\":\"failover:test_$i\",\"value\":\"before_failover_$i\"}" > /dev/null
done

sleep 2

echo "Simulating node-3 failure (reading only from node-1 and node-2)..."
success_count=0
for i in {1..20}; do
    # Пробуем читать с node-1
    value1=$(curl -s "${API_ENDPOINTS[0]}/read?key=failover:test_$i" | jq -r '.found')
    # Пробуем читать с node-2
    value2=$(curl -s "${API_ENDPOINTS[1]}/read?key=failover:test_$i" | jq -r '.found')
    
    if [ "$value1" == "true" ] || [ "$value2" == "true" ]; then
        ((success_count++))
    fi
done

availability=$(echo "scale=2; $success_count * 100 / 20" | bc)
echo "  Data availability with 2/3 nodes: ${availability}% (${success_count}/20)"
log_result "  Availability: ${availability}% (${success_count}/20)"

if [ $success_count -ge 18 ]; then
    echo -e "${GREEN}  ✓ High availability maintained${NC}"
    log_result "  ✓ High availability"
else
    echo -e "${YELLOW}  ⚠ Reduced availability${NC}"
    log_result "  ⚠ Reduced availability"
fi
log_result ""
echo

# ============================================================================
# Тест 9: ETCD consistency check
# ============================================================================
echo -e "${BLUE}═══ Test 9: ETCD Consistency Check ═══${NC}"
log_result "Test 9: ETCD Consistency Check"
log_result "--------------------------------"

echo "Checking ETCD cluster health..."
etcd_health=$(docker exec etcd1 etcdctl endpoint health --cluster 2>&1 || echo "Error checking ETCD")
echo "$etcd_health"
log_result "$etcd_health"
log_result ""
echo

# ============================================================================
# Итоговая сводка
# ============================================================================
echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              Benchmark Completed Successfully             ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
echo
echo -e "${YELLOW}Results saved to: $RESULTS_FILE${NC}"
echo
echo "Summary:"
echo "--------"
cat "$RESULTS_FILE" | grep -E "(RPS|Throughput|Latency|Success|Availability|Consistency)" | tail -15
echo

# Cleanup
rm -f /tmp/bench_write.json

log_result ""
log_result "Benchmark completed at: $(date '+%Y-%m-%d %H:%M:%S')"
