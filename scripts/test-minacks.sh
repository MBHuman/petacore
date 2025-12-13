#!/bin/bash

# Тест динамического изменения MIN_ACKS

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

API_URL="http://localhost:8081"

echo -e "${YELLOW}=== Тест динамического изменения MIN_ACKS ===${NC}\n"

# Тест 1: Проверяем текущий статус
echo -e "${YELLOW}Тест 1: Текущие настройки${NC}"
current_status=$(curl -s "$API_URL/status")
echo "$current_status" | jq '{node_id, total_nodes, min_acks, description}'
echo

# Тест 2: Устанавливаем MIN_ACKS=1 (слабая консистентность)
echo -e "${YELLOW}Тест 2: MIN_ACKS=1 (слабая консистентность)${NC}"
response=$(curl -s -X POST "$API_URL/config/min_acks" \
  -H "Content-Type: application/json" \
  -d '{"min_acks":1}')
echo "$response" | jq .

# Проверяем изменение
status=$(curl -s "$API_URL/status" | jq -r '.min_acks')
if [ "$status" == "1" ]; then
  echo -e "${GREEN}✓${NC} MIN_ACKS успешно изменен на 1"
else
  echo -e "${RED}✗${NC} Ошибка: MIN_ACKS = $status, ожидалось 1"
fi
echo

# Тест 3: Устанавливаем MIN_ACKS=0 (строгий quorum, по умолчанию)
echo -e "${YELLOW}Тест 3: MIN_ACKS=0 (строгий quorum)${NC}"
response=$(curl -s -X POST "$API_URL/config/min_acks" \
  -H "Content-Type: application/json" \
  -d '{"min_acks":0}')
echo "$response" | jq .

# Проверяем изменение (для 3 узлов должно быть 2)
status=$(curl -s "$API_URL/status" | jq -r '.min_acks')
total_nodes=$(curl -s "$API_URL/status" | jq -r '.total_nodes')
expected=$((total_nodes / 2 + 1))
if [ "$status" == "$expected" ]; then
  echo -e "${GREEN}✓${NC} MIN_ACKS успешно установлен на quorum: $status (N/2+1)"
else
  echo -e "${RED}✗${NC} Ошибка: MIN_ACKS = $status, ожидалось $expected"
fi
echo

# Тест 4: Устанавливаем MIN_ACKS=-1 (все узлы)
echo -e "${YELLOW}Тест 4: MIN_ACKS=-1 (все узлы)${NC}"
response=$(curl -s -X POST "$API_URL/config/min_acks" \
  -H "Content-Type: application/json" \
  -d '{"min_acks":-1}')
echo "$response" | jq .

# Проверяем изменение (должно быть равно total_nodes)
status=$(curl -s "$API_URL/status" | jq -r '.min_acks')
total_nodes=$(curl -s "$API_URL/status" | jq -r '.total_nodes')
if [ "$status" == "$total_nodes" ]; then
  echo -e "${GREEN}✓${NC} MIN_ACKS успешно установлен на все узлы: $status"
else
  echo -e "${RED}✗${NC} Ошибка: MIN_ACKS = $status, ожидалось $total_nodes"
fi
echo

# Тест 5: Тестируем чтение/запись с разными MIN_ACKS
echo -e "${YELLOW}Тест 5: Чтение/запись с разными MIN_ACKS${NC}"

# MIN_ACKS=1
echo -e "${BLUE}С MIN_ACKS=1:${NC}"
curl -s -X POST "$API_URL/config/min_acks" \
  -H "Content-Type: application/json" \
  -d '{"min_acks":1}' > /dev/null

curl -s -X POST "$API_URL/write" \
  -H "Content-Type: application/json" \
  -d '{"key":"test:minacks_1","value":"with_minacks_1"}' > /dev/null
echo -e "${GREEN}✓${NC} Запись выполнена с MIN_ACKS=1"

sleep 1

value=$(curl -s "$API_URL/read?key=test:minacks_1" | jq -r '.value')
echo -e "${GREEN}✓${NC} Чтение: value=$value"
echo

# MIN_ACKS=-1 (все узлы)
echo -e "${BLUE}С MIN_ACKS=-1 (все узлы):${NC}"
curl -s -X POST "$API_URL/config/min_acks" \
  -H "Content-Type: application/json" \
  -d '{"min_acks":-1}' > /dev/null

curl -s -X POST "$API_URL/write" \
  -H "Content-Type: application/json" \
  -d '{"key":"test:minacks_all","value":"with_all_nodes"}' > /dev/null
echo -e "${GREEN}✓${NC} Запись выполнена с MIN_ACKS=-1"

sleep 2

value=$(curl -s "$API_URL/read?key=test:minacks_all" | jq -r '.value')
echo -e "${GREEN}✓${NC} Чтение: value=$value"
echo

# Восстанавливаем значение по умолчанию
echo -e "${YELLOW}Восстановление значения по умолчанию...${NC}"
curl -s -X POST "$API_URL/config/min_acks" \
  -H "Content-Type: application/json" \
  -d '{"min_acks":0}' | jq .
echo

echo -e "${GREEN}=== Все тесты пройдены ===${NC}"
