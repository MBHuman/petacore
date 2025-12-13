# Уровни изоляции транзакций

Реализованы два уровня изоляции транзакций:

## Read Committed (по умолчанию)

При уровне изоляции **Read Committed**:
- Транзакция видит только committed (зафиксированные) данные
- При каждом чтении транзакция получает последнюю закоммиченную версию данных
- Возможен non-repeatable read: повторное чтение того же ключа в разных транзакциях может вернуть разные значения

### Пример использования

```go
// По умолчанию создается с Read Committed
s := storage.NewSimpleStorage()

s.RunTransaction(func(tx *core.Transaction) error {
    // Каждое чтение видит последние committed данные
    value, ok := tx.Read("key")
    
    tx.Write("key", "new_value")
    return nil
})
```

## Snapshot Isolation

При уровне изоляции **Snapshot Isolation**:
- Транзакция видит фиксированный snapshot данных на момент вызова `Begin()`
- Все чтения в рамках одной транзакции видят один и тот же snapshot
- Предотвращает non-repeatable read и phantom read

### Пример использования

```go
// Явно указываем Snapshot Isolation
s := storage.NewSimpleStorageWithIsolation(core.SnapshotIsolation)

s.RunTransaction(func(tx *core.Transaction) error {
    // Все чтения видят snapshot на момент начала транзакции
    value1, _ := tx.Read("key") // версия на момент Begin()
    
    // Даже если другая транзакция изменит "key" и закоммитит,
    // эта транзакция продолжит видеть старую версию
    value2, _ := tx.Read("key") // та же версия
    
    return nil
})
```

## Сравнение

| Характеристика | Read Committed | Snapshot Isolation |
|----------------|----------------|-------------------|
| Видимость данных | Последние committed | Фиксированный snapshot |
| Non-repeatable read | Возможен | Невозможен |
| Phantom read | Возможен | Невозможен |
| Производительность | Выше (нет блокировок) | Немного ниже |
| Использование памяти | Меньше | Больше (хранение snapshots) |

## Архитектура

Реализация основана на MVCC (Multi-Version Concurrency Control):

1. **LClock** - логические часы Лампорта для версионирования
2. **MVCC** - хранение множественных версий данных
3. **Transaction** - управление транзакциями с поддержкой уровней изоляции

### Ключевые отличия в реализации

**Read Committed:**
```go
func (tx *Transaction) Read(key string) (string, bool) {
    // Читаем с текущей версией (последние committed данные)
    currentVersion := int64(tx.logicalClock.Get())
    return tx.mvcc.Read(key, currentVersion)
}
```

**Snapshot Isolation:**
```go
func (tx *Transaction) Begin() {
    // Фиксируем версию snapshot
    version := tx.logicalClock.Get()
    tx.snapshotVersion = &version
}

func (tx *Transaction) Read(key string) (string, bool) {
    // Читаем с фиксированной версией snapshot
    return tx.mvcc.Read(key, int64(*tx.snapshotVersion))
}
```

## Тесты

Реализованы тесты, демонстрирующие поведение обоих уровней изоляции:

- `TestReadCommittedIsolation` - базовая проверка Read Committed
- `TestReadCommittedVsSnapshotIsolation` - сравнение уровней
- `TestReadCommittedNonRepeatableRead` - демонстрация non-repeatable read

Запуск тестов:
```bash
go test -v ./internal/storage -run TestReadCommitted
```
