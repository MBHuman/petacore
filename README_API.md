# PetaCore API Service

HTTP API сервис на основе PetaCore VClock распределенного хранилища.

## Быстрый старт

### 1. Сборка и запуск

```bash
# Собрать и запустить все сервисы
docker-compose -f docker-compose.api.yml up --build

# Или в фоновом режиме
docker-compose -f docker-compose.api.yml up --build -d
```

Это запустит:
- 3 узла ETCD (etcd1, etcd2, etcd3)
- 3 реплики API (api-node-1, api-node-2, api-node-3)

### Переменные окружения

- `NODE_ID` - идентификатор узла (обязательно)
- `ETCD_ENDPOINTS` - список адресов ETCD через запятую (по умолчанию: `localhost:2379`)
- `TOTAL_NODES` - общее количество узлов в кластере (по умолчанию: `3`)
- `MIN_ACKS` - минимальное количество подтверждений для quorum:
  - `0` или не указано - использовать `N/2 + 1` (строгий quorum, **по умолчанию**)
  - `-1` - использовать `N` (все узлы, максимальная консистентность)
  - `> 0` - конкретное значение
- `PORT` - порт HTTP сервера (по умолчанию: `8080`)

**Примеры настройки MIN_ACKS:**

```bash
# Строгий quorum (большинство) - по умолчанию
MIN_ACKS=0  # или не указывать - будет N/2 + 1
# Для 3 узлов: min_acks = 2

# Более слабая консистентность (быстрее, но менее надежно)
MIN_ACKS=1  # достаточно одного подтверждения

# Все узлы (максимальная консистентность)
MIN_ACKS=-1  # требуется подтверждение от всех узлов
# Для 3 узлов: min_acks = 3
```

### 2. Проверка здоровья

```bash
# Проверить все узлы
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

# Статус узлов
curl http://localhost:8081/status | jq
```

## API Endpoints

### POST /write
Записать значение

```bash
curl -X POST http://localhost:8081/write \
  -H "Content-Type: application/json" \
  -d '{"key":"user:123","value":"Alice"}'
```

### GET /read
Прочитать значение с quorum проверкой

```bash
curl "http://localhost:8081/read?key=user:123" | jq
```

Ответ:
```json
{
  "key": "user:123",
  "value": "Alice",
  "found": true
}
```

### GET /status
Статус узла

```bash
curl http://localhost:8081/status | jq
```

Ответ:
```json
{
  "node_id": "api-node-1",
  "is_synced": true,
  "total_nodes": 3,
  "min_acks": 2,
  "description": "VClock distributed storage with quorum (min 2/3 nodes)"
}
```

### GET /health
Health check для load balancer

```bash
curl http://localhost:8081/health
```

### POST/PUT /config/min_acks
Динамическое изменение минимального количества подтверждений (quorum)

```bash
# Установить строгий quorum (N/2 + 1)
curl -X POST http://localhost:8081/config/min_acks \
  -H "Content-Type: application/json" \
  -d '{"min_acks":0}'

# Установить все узлы (максимальная консистентность)
curl -X POST http://localhost:8081/config/min_acks \
  -H "Content-Type: application/json" \
  -d '{"min_acks":-1}'

# Установить конкретное значение
curl -X POST http://localhost:8081/config/min_acks \
  -H "Content-Type: application/json" \
  -d '{"min_acks":2}'
```

Ответ:
```json
{
  "status": "ok",
  "old_min_acks": 2,
  "new_min_acks": 3,
  "total_nodes": 3,
  "description": "Set to all nodes (3/3)"
}
```

**Значения min_acks:**
- `0` - использовать N/2 + 1 (строгий quorum, по умолчанию)
- `-1` - использовать N (все узлы, максимальная консистентность)
- `> 0` - конкретное значение

## Тестирование

### Автоматические тесты

```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

Тесты проверяют:
1. Health check всех узлов
2. Статус синхронизации
3. Запись на один узел
4. Чтение со всех узлов (проверка репликации)
5. Множественные записи
6. Параллельные записи
7. Консистентность после конкурентных операций

### Бенчмарки

```bash
chmod +x scripts/benchmark-api.sh
./scripts/benchmark-api.sh
```

Бенчмарки измеряют:
1. **Read Performance** - производительность чтения (ожидается ~100k+ req/s)
2. **Write Performance** - производительность записи (~200-500 req/s)
3. **Distributed Load** - распределение нагрузки по узлам
4. **Latency Distribution** - перцентили задержек (p50, p90, p99)
5. **Maximum Throughput** - максимальная пропускная способность
6. **Mixed Workload** - смешанная нагрузка (read + write)

## Ручное тестирование

### Сценарий 1: Базовая репликация

```bash
# Запишите на node-1
curl -X POST http://localhost:8081/write \
  -H "Content-Type: application/json" \
  -d '{"key":"test:replication","value":"from_node_1"}'

# Подождите 1-2 секунды

# Прочитайте с node-2
curl "http://localhost:8082/read?key=test:replication" | jq

# Прочитайте с node-3
curl "http://localhost:8083/read?key=test:replication" | jq
```

**Ожидаемый результат:** Все узлы возвращают одно и то же значение.

### Сценарий 2: Quorum чтение

```bash
# Запишите данные
curl -X POST http://localhost:8081/write \
  -H "Content-Type: application/json" \
  -d '{"key":"test:quorum","value":"quorum_test"}'

# Сразу прочитайте (может не быть quorum)
curl "http://localhost:8082/read?key=test:quorum" | jq

# Подождите 2 секунды и прочитайте снова
sleep 2
curl "http://localhost:8082/read?key=test:quorum" | jq
```

**Ожидаемый результат:** После синхронизации все узлы видят данные.

### Сценарий 3: Конкурентные записи

```bash
# Запустите 10 параллельных записей
for i in {1..10}; do
  curl -X POST http://localhost:808$((1 + RANDOM % 3))/write \
    -H "Content-Type: application/json" \
    -d "{\"key\":\"counter\",\"value\":\"$i\"}" &
done
wait

# Подождите синхронизации
sleep 3

# Прочитайте финальное значение
curl "http://localhost:8081/read?key=counter" | jq
curl "http://localhost:8082/read?key=counter" | jq
curl "http://localhost:8083/read?key=counter" | jq
```

**Ожидаемый результат:** Все узлы возвращают одно и то же (последнее) значение.

### Сценарий 4: Отказоустойчивость

```bash
# Остановите один узел
docker stop petacore-api-3

# Запишите данные на оставшиеся узлы
curl -X POST http://localhost:8081/write \
  -H "Content-Type: application/json" \
  -d '{"key":"test:failover","value":"with_2_nodes"}'

# Прочитайте с node-2
curl "http://localhost:8082/read?key=test:failover" | jq

# Запустите node-3 обратно
docker start petacore-api-3

# Подождите синхронизации
sleep 5

# Прочитайте с восстановленного node-3
curl "http://localhost:8083/read?key=test:failover" | jq
```

**Ожидаемый результат:** Система работает с 2 из 3 узлов, восстановленный узел синхронизируется.

### Сценарий 5: Настройка MIN_ACKS

MIN_ACKS позволяет настроить уровень консистентности:

**Запуск с MIN_ACKS=1 (слабая консистентность, высокая скорость):**

```bash
# В docker-compose.api.yml добавьте:
# environment:
#   - MIN_ACKS=1

# Или запустите локально
MIN_ACKS=1 NODE_ID=test-node PORT=8080 go run cmd/api/main.go
```

**Проверка работы:**

```bash
# Запишите данные
curl -X POST http://localhost:8080/write \
  -H "Content-Type: application/json" \
  -d '{"key":"test:minacks","value":"with_minacks_1"}'

# Немедленно прочитайте (будет доступно сразу)
curl "http://localhost:8080/read?key=test:minacks" | jq

# Проверьте статус
curl http://localhost:8080/status | jq
# Ожидаемый min_acks: 1
```

**Запуск с MIN_ACKS=3 (строгая консистентность):**

```bash
MIN_ACKS=3 TOTAL_NODES=3 NODE_ID=strict-node PORT=8090 go run cmd/api/main.go
```

**Сравнение производительности:**

| MIN_ACKS | Консистентность | Скорость чтения | Доступность | Рекомендуется для |
|----------|----------------|-----------------|-------------|-------------------|
| 1        | Слабая         | Очень высокая   | Высокая     | Кеши, метрики     |
| 2 (N/2+1)| Строгая        | Высокая         | Средняя     | **По умолчанию**  |
| N (все)  | Очень строгая  | Средняя         | Низкая      | Критичные данные  |

## Мониторинг

### Логи

```bash
# Все логи
docker-compose -f docker-compose.api.yml logs -f

# Конкретный узел
docker-compose -f docker-compose.api.yml logs -f api-node-1

# ETCD логи
docker-compose -f docker-compose.api.yml logs -f etcd1
```

### Метрики

```bash
# Статус всех узлов
for port in 8081 8082 8083; do
  echo "Node on port $port:"
  curl -s "http://localhost:$port/status" | jq
done
```

### Данные в ETCD

```bash
# Посмотреть все ключи
docker exec -it etcd1 etcdctl get --prefix "petacore-api" --keys-only

# Посмотреть конкретный ключ
docker exec -it etcd1 etcdctl get --prefix "petacore-api/test:replication" --print-value-only | jq
```

## Производительность

### Ожидаемые показатели

- **Read Latency**: 0.5-2ms (локальный кеш с quorum проверкой)
- **Write Latency**: 5-10ms (ETCD + синхронизация)
- **Read Throughput**: 50,000-100,000+ req/s на узел
- **Write Throughput**: 200-500 req/s на узел
- **Quorum**: 2 из 3 узлов (N/2 + 1)

### Факторы производительности

**Чтение (быстро):**
- Локальный MVCC кеш
- Quorum проверка без сетевых вызовов
- Vector Clock хранится с данными

**Запись (медленно):**
- Синхронная запись в ETCD
- Raft консенсус в ETCD кластере
- Асинхронная репликация на другие узлы

## Остановка

```bash
# Остановить все сервисы
docker-compose -f docker-compose.api.yml down

# Удалить volumes (очистить данные)
docker-compose -f docker-compose.api.yml down -v
```

## Troubleshooting

### Узел не синхронизируется

```bash
# Проверьте ETCD
docker exec -it etcd1 etcdctl endpoint health

# Перезапустите узел
docker-compose -f docker-compose.api.yml restart api-node-1
```

### Чтения не находят данные

1. Проверьте, что запись прошла успешно
2. Подождите 1-2 секунды для синхронизации
3. Проверьте, что quorum достигнут (minAcks = 2 для 3 узлов)
4. Посмотрите логи узла

### Медленные записи

Это нормально! Записи идут через ETCD с Raft консенсусом. Для увеличения производительности записей используйте батчинг или асинхронную запись (если допустимо eventual consistency).

## Архитектура

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  api-node-1 │  │  api-node-2 │  │  api-node-3 │
│   (8081)    │  │   (8082)    │  │   (8083)    │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
         ┌──────────────┴──────────────┐
         │         ETCD Cluster        │
         │  (etcd1, etcd2, etcd3)     │
         └─────────────────────────────┘

- Каждый API узел имеет локальный MVCC кеш
- Записи идут через ETCD (single source of truth)
- Чтения из локального кеша с quorum проверкой
- Vector Clock синхронизируется через watch
```
