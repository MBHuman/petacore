# Перегрузка функций в SDK

Petacore SDK поддерживает перегрузку функций (function overloading) — возможность иметь несколько реализаций одной функции с разными типами аргументов.

## Как это работает

### Регистрация перегруженных функций

Когда вы регистрируете функции с одинаковым именем но разными OID и типами аргументов, они сохраняются как разные реализации:

```go
func RegisterFunctions(registry *psdk.FunctionRegistry) error {
    funcs := []psdk.IFunction{
        &MaxFunction{},     // MAX(float8) -> float8
        &MaxFunctionInt{},  // MAX(int4) -> int4
        &MaxFunctionText{}, // MAX(text) -> text
    }
    
    for _, fn := range funcs {
        if err := registry.Register(fn); err != nil {
            return err
        }
    }
    return nil
}
```

### Пример реализации перегруженных функций

```go
// MAX для FLOAT8
type MaxFunction struct {
    *psdk.BaseFunction
}

func (f *MaxFunction) GetFunction() *psdk.Function {
    return &psdk.Function{
        OID:         2116,
        ProName:     "MAX",
        IsAggregate: true,
        ProArgTypes: []psdk.OID{psdk.PTypeFloat8},
        ProRetType:  psdk.PTypeFloat8,
    }
}

// MAX для INT4
type MaxFunctionInt struct {
    *psdk.BaseFunction
}

func (f *MaxFunctionInt) GetFunction() *psdk.Function {
    return &psdk.Function{
        OID:         2117,
        ProName:     "MAX",
        IsAggregate: true,
        ProArgTypes: []psdk.OID{psdk.PTypeInt4},
        ProRetType:  psdk.PTypeInt4,
    }
}

// MAX для TEXT
type MaxFunctionText struct {
    *psdk.BaseFunction
}

func (f *MaxFunctionText) GetFunction() *psdk.Function {
    return &psdk.Function{
        OID:         2129,
        ProName:     "MAX",
        IsAggregate: true,
        ProArgTypes: []psdk.OID{psdk.PTypeText},
        ProRetType:  psdk.PTypeText,
    }
}
```

## API для работы с перегруженными функциями

### GetByName(name string)

Возвращает **первую** зарегистрированную функцию с указанным именем:

```go
fn, exists := registry.GetByName("MAX")
// Вернет MaxFunction (float8 версию)
```

### GetByNameAndArgTypes(name string, argTypes []OID)

Возвращает функцию с **точным совпадением** типов аргументов:

```go
// Получить MAX для int4
fn, exists := registry.GetByNameAndArgTypes("MAX", []psdk.OID{psdk.PTypeInt4})

// Получить MAX для text
fn, exists := registry.GetByNameAndArgTypes("MAX", []psdk.OID{psdk.PTypeText})
```

### GetAllByName(name string)

Возвращает **все** перегрузки функции:

```go
allMaxFuncs := registry.GetAllByName("MAX")
// Вернет все три версии: MaxFunction, MaxFunctionInt, MaxFunctionText
```

## Важные моменты

1. **Уникальные OID**: Каждая перегрузка должна иметь уникальный OID
2. **Типы аргументов**: Функции различаются по типам аргументов (ProArgTypes)
3. **Порядок регистрации**: При использовании GetByName возвращается первая зарегистрированная функция
4. **Автоматический выбор**: В будущем система сможет автоматически выбирать подходящую перегрузку на основе типов аргументов в запросе

## Пример использования

```sql
-- Для разных типов данных будут использоваться разные реализации MAX:

SELECT MAX(price) FROM products;        -- MAX(float8)
SELECT MAX(id) FROM users;              -- MAX(int4) 
SELECT MAX(name) FROM customers;        -- MAX(text)
```

## Совместимость с PostgreSQL

Система перегрузки повторяет подход PostgreSQL, где функции идентифицируются по:
- Имени функции
- Количеству аргументов
- Типам аргументов

Это позволяет создавать SQL-совместимые реализации функций.
