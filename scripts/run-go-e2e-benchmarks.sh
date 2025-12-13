#!/bin/bash

# Скрипт для запуска Go E2E бенчмарков VClock с реальным ETCD

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║     PetaCore VClock E2E Go Benchmarks                     ║${NC}"
echo -e "${YELLOW}╚═══════════════════════════════════════════════════════════╝${NC}"
echo

# Проверка доступности ETCD
echo -e "${BLUE}Checking ETCD cluster...${NC}"
if ! docker ps | grep -q etcd1; then
    echo -e "${RED}Error: ETCD cluster is not running${NC}"
    echo "Start it with: docker-compose -f docker-compose.api.yml up -d etcd1 etcd2 etcd3"
    exit 1
fi

# Проверка здоровья ETCD
if ! docker exec etcd1 etcdctl endpoint health > /dev/null 2>&1; then
    echo -e "${RED}Error: ETCD is not healthy${NC}"
    exit 1
fi

echo -e "${GREEN}✓ ETCD cluster is healthy${NC}"
echo

# Создаем директорию для результатов
RESULTS_DIR="benchmark-results"
mkdir -p "$RESULTS_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="$RESULTS_DIR/go_e2e_vclock_$TIMESTAMP.txt"

echo -e "${YELLOW}Running benchmarks...${NC}"
echo "Results will be saved to: $RESULTS_FILE"
echo

# Функция для запуска бенчмарка
run_bench() {
    local name="$1"
    local pattern="$2"
    local benchtime="${3:-10s}"
    
    echo -e "${BLUE}═══ $name ═══${NC}"
    echo "Pattern: $pattern, Time: $benchtime"
    echo
    
    go test -bench="$pattern" \
        -benchmem \
        -benchtime="$benchtime" \
        -timeout=30m \
        ./internal/storage/ \
        2>&1 | tee -a "$RESULTS_FILE"
    
    echo
    echo "---"
    echo
}

# Основные бенчмарки
cat > "$RESULTS_FILE" <<EOF
========================================
PetaCore VClock E2E Go Benchmarks
========================================
Date: $(date '+%Y-%m-%d %H:%M:%S')

EOF

echo -e "${YELLOW}1. Write Performance Tests${NC}"
run_bench "Write with Quorum" "BenchmarkE2E_WriteQuorum$" "10s"
run_bench "Write All Nodes" "BenchmarkE2E_WriteAllNodes$" "10s"
run_bench "Write Weak Consistency" "BenchmarkE2E_WriteWeakConsistency$" "10s"

echo -e "${YELLOW}2. Read Performance Tests${NC}"
run_bench "Read with Quorum" "BenchmarkE2E_ReadWithQuorum$" "10s"

echo -e "${YELLOW}3. Mixed Workload Tests${NC}"
run_bench "Mixed Workload (70% read, 30% write)" "BenchmarkE2E_MixedWorkload$" "10s"

echo -e "${YELLOW}4. Distributed Tests${NC}"
run_bench "Distributed Writes (3 nodes)" "BenchmarkE2E_DistributedWrites$" "10s"

echo -e "${YELLOW}5. Vector Clock Synchronization${NC}"
run_bench "VClock Sync Speed" "BenchmarkE2E_VClockSync$" "5s"

echo -e "${YELLOW}6. Concurrency Tests${NC}"
run_bench "Concurrent Writes" "BenchmarkE2E_ConcurrentWrites" "5s"

echo -e "${YELLOW}7. Batch Tests${NC}"
run_bench "Batch Writes" "BenchmarkE2E_BatchWrites" "5s"

echo -e "${YELLOW}8. Consistency Tests${NC}"
run_bench "Consistency Check" "BenchmarkE2E_ConsistencyCheck$" "5s"

echo -e "${YELLOW}9. High Contention Tests${NC}"
run_bench "High Contention Write" "BenchmarkE2E_HighContentionWrite$" "5s"

echo -e "${YELLOW}10. Performance Metrics${NC}"
run_bench "Latency Measurement" "BenchmarkE2E_Latency$" "10s"
run_bench "Throughput Measurement" "BenchmarkE2E_Throughput$" "10s"

# Создаем сводку
echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              Benchmarks Completed                         ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
echo
echo -e "${YELLOW}Results saved to: $RESULTS_FILE${NC}"
echo
echo "Summary of key metrics:"
echo "----------------------"
grep -E "(BenchmarkE2E|ns/op|ops/sec|µs/op|B/op|allocs/op)" "$RESULTS_FILE" | \
    grep -v "^#" | \
    awk '{
        if ($1 ~ /^BenchmarkE2E/) {
            bench=$1
        } else if ($2 ~ /ns\/op/) {
            printf "%-50s %10s ns/op\n", bench, $1
        } else if ($2 ~ /ops\/sec/) {
            printf "%-50s %10s ops/sec\n", bench, $1
        }
    }' | tail -20

echo
echo "For full results, see: $RESULTS_FILE"
echo

# Опционально: генерируем графики если доступен benchstat
if command -v benchstat &> /dev/null; then
    echo -e "${BLUE}Generating comparison report...${NC}"
    benchstat "$RESULTS_FILE" 2>/dev/null || true
fi
