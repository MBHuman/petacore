# Query Planner Architecture

Новая архитектура выполнения запросов основана на разделении планирования и выполнения.

## Компоненты

### 1. Planner (plan ировщик)
- **Файл**: `internal/runtime/planner/planner.go`
- **Функция**: `CreateQueryPlan(stmt *statements.SelectStatement, ctx PlannerContext) (*QueryPlan, error)`
- Преобразует AST (абстрактное синтаксическое дерево) в план выполнения
- Оптимизирует порядок операций
- Создает дерево узлов плана

### 2. Plan Nodes (узлы плана)
- **Файл**: `internal/runtime/planner/plan.go`

#### Типы узлов:
- **ScanPlanNode**: Сканирование таблицы
- **ValuesPlanNode**: SELECT без таблицы (SELECT 1, 2, 3)
- **ProjectPlanNode**: Выбор колонок (SELECT a, b)
- **FilterPlanNode**: Фильтрация (WHERE)
- **JoinPlanNode**: Соединение таблиц (JOIN)
- **AggregatePlanNode**: Агрегация (GROUP BY)
- **SortPlanNode**: Сортировка (ORDER BY)
- **LimitPlanNode**: Ограничение (LIMIT/OFFSET)
- **UnionPlanNode**: Объединение (UNION)
- **IntersectPlanNode**: Пересечение (INTERSECT)
- **ExceptPlanNode**: Разность (EXCEPT)
- **SubqueryPlanNode**: Подзапрос в FROM

### 3. Executor (исполнитель)
- **Файл**: `internal/runtime/planner/executor.go`
- **Функция**: `ExecutePlan(plan *QueryPlan, ctx ExecutorContext) (*table.ExecuteResult, error)`
- Рекурсивно выполняет узлы плана
- Каждый узел имеет свою функцию выполнения: `executeScan`, `executeProject`, и т.д.

### 4. Set Operations (операции множеств)
- **Файл**: `internal/runtime/planner/set_operations.go`
- Реализует UNION, INTERSECT, EXCEPT
- Поддерживает версии с ALL
- Удаляет дубликаты для обычных версий

## Преимущества новой архитектуры

1. **Модульность**: Каждая операция — отдельный узел плана
2. **Переиспользование**: Базовые операции используются повторно
3. **Тестируемость**: Можно тестировать каждый узел отдельно
4. **Расширяемость**: Легко добавлять новые операции
5. **Оптимизация**: В будущем можно добавить оптимизатор планов

## Пример использования

```go
// 1. Парсинг SQL
stmt, err := visitor.ParseSQL("SELECT a, b FROM table1 UNION SELECT c, d FROM table2")

// 2. Создание плана
plannerCtx := planner.PlannerContext{
    Database: "mydb",
    Schema:   "public",
}
plan, err := planner.CreateQueryPlan(stmt.(*statements.SelectStatement), plannerCtx)

// 3. Выполнение плана
executorCtx := planner.ExecutorContext{
    Database: "mydb",
    Schema:   "public",
    Storage:  storage,
}
result, err := planner.ExecutePlan(plan, executorCtx)
```

## Структура плана

Для запроса `SELECT a, b FROM t1 WHERE x > 10 ORDER BY a LIMIT 5`:

```
Limit
  Sort
    Project
      Filter
        Scan
```

Для запроса `SELECT * FROM t1 UNION SELECT * FROM t2`:

```
Union
  Project
    Scan (t1)
  Project
    Scan (t2)
```

## Миграция со старого кода

Старые функции `ExecuteSelect`, `ExecuteNormalTable`, `ExecuteSelectWithoutTable` более не используются.
Весь код теперь работает через `ExecuteSelectWithPlanner` в `executor/executor.go`.

## TODO

- [ ] Реализовать полноценную агрегацию в `processGroupByAndAggregates`
- [ ] Добавить поддержку ON условий в JOIN
- [ ] Реализовать RIGHT и FULL JOIN
- [ ] Добавить оптимизатор планов
- [ ] Добавить статистику выполнения
- [ ] Добавить EXPLAIN для просмотра планов
