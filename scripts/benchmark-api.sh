#!/bin/bash

# Бенчмарк для PetaCore API

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

API_ENDPOINTS=(
    "http://localhost:8081"
    "http://localhost:8082"
    "http://localhost:8083"
)

echo -e "${YELLOW}=== PetaCore API Benchmark ===${NC}\n"

# Проверка наличия apache bench
if ! command -v ab &> /dev/null; then
    echo "Installing apache bench..."
    sudo apt-get update && sudo apt-get install -y apache2-utils
fi

# Бенчмарк 1: Чтение (GET)
echo -e "${YELLOW}Benchmark 1: Read Performance${NC}"
echo "Warming up..."
for i in {1..100}; do
    curl -s "${API_ENDPOINTS[0]}/read?key=bench:test" > /dev/null
done

echo "Running benchmark: 1000 requests, concurrency 10"
ab -n 1000 -c 10 -q "${API_ENDPOINTS[0]}/read?key=bench:test" 2>&1 | grep -E "(Requests per second|Time per request|Transfer rate)"
echo ""

# Бенчмарк 2: Запись (POST)
echo -e "${YELLOW}Benchmark 2: Write Performance${NC}"
echo "Creating test data file..."
echo '{"key":"bench:write","value":"benchmark test value"}' > /tmp/petacore-bench.json

echo "Running benchmark: 500 requests, concurrency 5"
ab -n 500 -c 5 -p /tmp/petacore-bench.json -T "application/json" -q "${API_ENDPOINTS[0]}/write" 2>&1 | grep -E "(Requests per second|Time per request|Transfer rate)"
echo ""

# Бенчмарк 3: Распределенная нагрузка
echo -e "${YELLOW}Benchmark 3: Distributed Load${NC}"
echo "Testing load distribution across 3 nodes..."

for i in "${!API_ENDPOINTS[@]}"; do
    endpoint="${API_ENDPOINTS[$i]}"
    echo "Node $((i+1)):"
    ab -n 300 -c 10 -q "$endpoint/read?key=bench:distributed" 2>&1 | grep "Requests per second"
done
echo ""

# Бенчмарк 4: Latency test
echo -e "${YELLOW}Benchmark 4: Latency Distribution${NC}"
echo "Testing latency percentiles..."

ab -n 1000 -c 10 -q -g /tmp/petacore-latency.tsv "${API_ENDPOINTS[0]}/read?key=bench:latency" 2>&1 | grep -E "(50%|90%|99%|100%)"
echo ""

# Бенчмарк 5: Throughput под нагрузкой
echo -e "${YELLOW}Benchmark 5: Maximum Throughput${NC}"
echo "Testing maximum throughput with high concurrency..."

echo "Concurrency 50:"
ab -n 1000 -c 50 -q "${API_ENDPOINTS[0]}/read?key=bench:throughput" 2>&1 | grep "Requests per second"

echo "Concurrency 100:"
ab -n 1000 -c 100 -q "${API_ENDPOINTS[0]}/read?key=bench:throughput" 2>&1 | grep "Requests per second"
echo ""

# Бенчмарк 6: Mixed workload
echo -e "${YELLOW}Benchmark 6: Mixed Read/Write Workload${NC}"
echo "Simulating mixed workload..."

# Параллельно запускаем чтения и записи
(
    for i in {1..500}; do
        curl -s "${API_ENDPOINTS[0]}/read?key=bench:mixed$((i % 10))" > /dev/null
    done
) &

(
    for i in {1..100}; do
        curl -s -X POST "${API_ENDPOINTS[1]}/write" \
            -H "Content-Type: application/json" \
            -d "{\"key\":\"bench:mixed$((i % 10))\",\"value\":\"value$i\"}" > /dev/null
    done
) &

wait

echo -e "${GREEN}✓${NC} Mixed workload completed (500 reads + 100 writes)"
echo ""

# Очистка
rm -f /tmp/petacore-bench.json /tmp/petacore-latency.tsv

echo -e "${GREEN}=== Benchmark completed ===${NC}"
echo ""
echo "Summary:"
echo "- Read performance: High throughput expected (local cache)"
echo "- Write performance: Limited by ETCD + Vector Clock sync"
echo "- Distributed: All nodes should show similar read performance"
