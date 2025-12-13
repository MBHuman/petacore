#!/bin/bash

# Скрипт для запуска бенчмарков распределённой СУБД PetaCore

set -e

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== PetaCore Distributed Storage Benchmarks ===${NC}\n"

# Проверяем, что ETCD запущен
echo -e "${YELLOW}Проверка ETCD...${NC}"
if ! docker-compose ps | grep -q etcd1.*Up; then
    echo "ETCD не запущен. Запускаем..."
    docker-compose up -d
    echo "Ожидание запуска ETCD..."
    sleep 5
fi
echo -e "${GREEN}✓ ETCD запущен${NC}\n"

# Переходим в директорию с тестами
cd "$(dirname "$0")/internal/storage"

# Функция для запуска бенчмарка
run_benchmark() {
    local name=$1
    local pattern=$2
    echo -e "${BLUE}### $name ###${NC}"
    go test -bench="$pattern" -benchmem -benchtime=3s -timeout=30m
    echo ""
}

# Меню выбора
if [ "$1" == "all" ]; then
    echo -e "${YELLOW}Запуск всех бенчмарков...${NC}\n"
    run_benchmark "Все бенчмарки" "^BenchmarkDistributed"
    
elif [ "$1" == "quick" ]; then
    echo -e "${YELLOW}Быстрые бенчмарки...${NC}\n"
    run_benchmark "Запись" "^BenchmarkDistributedWrite$"
    run_benchmark "Чтение" "^BenchmarkDistributedRead$"
    
elif [ "$1" == "write" ]; then
    run_benchmark "Запись" "^BenchmarkDistributedWrite$"
    
elif [ "$1" == "read" ]; then
    run_benchmark "Чтение" "^BenchmarkDistributedRead$"
    
elif [ "$1" == "readwrite" ]; then
    run_benchmark "Чтение-Запись" "^BenchmarkDistributedReadWrite$"
    
elif [ "$1" == "transaction" ]; then
    run_benchmark "Транзакции" "^BenchmarkDistributedTransaction$"
    
elif [ "$1" == "concurrent" ]; then
    run_benchmark "Конкурентные записи" "^BenchmarkDistributedConcurrentWrites$"
    run_benchmark "Конкурентные чтения" "^BenchmarkDistributedConcurrentReads$"
    
elif [ "$1" == "hotkey" ]; then
    run_benchmark "Hot Key" "^BenchmarkDistributedHotKey$"
    
elif [ "$1" == "isolation" ]; then
    run_benchmark "Уровни изоляции" "^BenchmarkDistributedIsolationLevels"
    
elif [ "$1" == "multinode" ]; then
    run_benchmark "Несколько узлов" "^BenchmarkDistributedMultiNode$"
    
elif [ "$1" == "batch" ]; then
    run_benchmark "Пакетные записи" "^BenchmarkDistributedBatchWrites"
    
elif [ "$1" == "latency" ]; then
    run_benchmark "Задержки" "^BenchmarkDistributedLatency$"
    
elif [ "$1" == "compare" ]; then
    echo -e "${YELLOW}Сравнение локальной и распределённой СУБД...${NC}\n"
    echo -e "${BLUE}### Локальная СУБД ###${NC}"
    go test -bench="^BenchmarkWrite$" -benchmem -benchtime=3s
    echo ""
    echo -e "${BLUE}### Распределённая СУБД ###${NC}"
    go test -bench="^BenchmarkDistributedWrite$" -benchmem -benchtime=3s
    
else
    echo -e "Использование: $0 <command>"
    echo ""
    echo "Команды:"
    echo "  all         - Запустить все бенчмарки"
    echo "  quick       - Быстрые бенчмарки (запись + чтение)"
    echo "  write       - Производительность записи"
    echo "  read        - Производительность чтения"
    echo "  readwrite   - Смешанные операции"
    echo "  transaction - Производительность транзакций"
    echo "  concurrent  - Конкурентный доступ"
    echo "  hotkey      - Работа с hot key"
    echo "  isolation   - Сравнение уровней изоляции"
    echo "  multinode   - Несколько узлов"
    echo "  batch       - Пакетные операции"
    echo "  latency     - Измерение задержек"
    echo "  compare     - Сравнение с локальной СУБД"
    echo ""
    echo "Примеры:"
    echo "  $0 quick"
    echo "  $0 concurrent"
    echo "  $0 compare"
    exit 1
fi

echo -e "${GREEN}=== Бенчмарки завершены ===${NC}"
