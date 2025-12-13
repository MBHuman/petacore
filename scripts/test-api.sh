#!/bin/bash

# Тестовый скрипт для PetaCore API

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

API_ENDPOINTS=(
    "http://localhost:8081"
    "http://localhost:8082"
    "http://localhost:8083"
)

echo -e "${YELLOW}=== PetaCore API Test ===${NC}\n"

# Тест 1: Проверка здоровья всех узлов
echo -e "${YELLOW}Test 1: Health check${NC}"
for endpoint in "${API_ENDPOINTS[@]}"; do
    response=$(curl -s -o /dev/null -w "%{http_code}" "$endpoint/health")
    if [ "$response" == "200" ]; then
        echo -e "${GREEN}✓${NC} $endpoint is healthy"
    else
        echo -e "${RED}✗${NC} $endpoint is not healthy (HTTP $response)"
    fi
done
echo ""

# Тест 2: Проверка статуса
echo -e "${YELLOW}Test 2: Status check${NC}"
for endpoint in "${API_ENDPOINTS[@]}"; do
    status=$(curl -s "$endpoint/status" | jq -r '.node_id + " synced=" + (.is_synced|tostring)')
    echo -e "${GREEN}✓${NC} $status"
done
echo ""

# Тест 3: Запись на первый узел
echo -e "${YELLOW}Test 3: Write to node-1${NC}"
write_response=$(curl -s -X POST "${API_ENDPOINTS[0]}/write" \
    -H "Content-Type: application/json" \
    -d '{"key":"test:key1","value":"hello from node-1"}')
write_status=$(echo "$write_response" | jq -r '.status')
if [ "$write_status" == "ok" ]; then
    echo -e "${GREEN}✓${NC} Write successful"
else
    echo -e "${RED}✗${NC} Write failed"
    echo "$write_response" | jq .
fi
echo ""

# Ждем синхронизации
echo -e "${YELLOW}Waiting for sync...${NC}"
sleep 2

# Тест 4: Чтение со всех узлов
echo -e "${YELLOW}Test 4: Read from all nodes${NC}"
for i in "${!API_ENDPOINTS[@]}"; do
    endpoint="${API_ENDPOINTS[$i]}"
    response=$(curl -s "$endpoint/read?key=test:key1")
    value=$(echo "$response" | jq -r '.value')
    found=$(echo "$response" | jq -r '.found')
    
    if [ "$found" == "true" ]; then
        echo -e "${GREEN}✓${NC} Node $((i+1)): value=$value"
    else
        echo -e "${RED}✗${NC} Node $((i+1)): key not found"
    fi
done
echo ""

# Тест 5: Множественные записи
echo -e "${YELLOW}Test 5: Multiple writes${NC}"
for i in {1..5}; do
    node_idx=$((i % 3))
    endpoint="${API_ENDPOINTS[$node_idx]}"
    curl -s -X POST "$endpoint/write" \
        -H "Content-Type: application/json" \
        -d "{\"key\":\"test:multi$i\",\"value\":\"value$i\"}" > /dev/null
    echo -e "${GREEN}✓${NC} Written test:multi$i to node-$((node_idx+1))"
done
echo ""

sleep 2

# Тест 6: Чтение множественных записей
echo -e "${YELLOW}Test 6: Read multiple keys${NC}"
for i in {1..5}; do
    response=$(curl -s "${API_ENDPOINTS[0]}/read?key=test:multi$i")
    value=$(echo "$response" | jq -r '.value')
    found=$(echo "$response" | jq -r '.found')
    
    if [ "$found" == "true" ] && [ "$value" == "value$i" ]; then
        echo -e "${GREEN}✓${NC} test:multi$i = $value"
    else
        echo -e "${RED}✗${NC} test:multi$i failed"
    fi
done
echo ""

# Тест 7: Параллельные записи
echo -e "${YELLOW}Test 7: Concurrent writes${NC}"
for i in {1..10}; do
    node_idx=$((RANDOM % 3))
    endpoint="${API_ENDPOINTS[$node_idx]}"
    curl -s -X POST "$endpoint/write" \
        -H "Content-Type: application/json" \
        -d "{\"key\":\"test:concurrent\",\"value\":\"iteration$i\"}" > /dev/null &
done
wait
echo -e "${GREEN}✓${NC} 10 concurrent writes completed"
echo ""

sleep 2

# Тест 8: Чтение после параллельных записей
echo -e "${YELLOW}Test 8: Read after concurrent writes${NC}"
response=$(curl -s "${API_ENDPOINTS[1]}/read?key=test:concurrent")
value=$(echo "$response" | jq -r '.value')
found=$(echo "$response" | jq -r '.found')

if [ "$found" == "true" ]; then
    echo -e "${GREEN}✓${NC} Final value: $value"
else
    echo -e "${RED}✗${NC} Key not found after concurrent writes"
fi
echo ""

echo -e "${GREEN}=== All tests completed ===${NC}"
