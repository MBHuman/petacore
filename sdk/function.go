package psdk

import (
	"context"
	"fmt"
	"sync"
)

type BaseFunction struct{}

func (bf *BaseFunction) GetAggFunc() (IAggFunction, error) {
	return nil, nil
}

type IFunction interface {
	Execute(ctx context.Context, args ...any) (any, error)
	GetFunction() *Function
	GetAggFunc() (IAggFunction, error)
}

type IAggFunction interface {
	Execute(ctx context.Context, aggVal interface{}, args ...any) (any, error)
}

// Function represents a user-defined function
type Function struct {
	OID         OID    // Unique identifier for the function
	ProName     string // Function name
	ProArgTypes []OID  // Argument types (as OIDs)
	ProRetType  OID    // Return type (as OID)
	IsAggregate bool

	Meta FunctionMeta // Additional metadata about the function
}

type FunctionAgg struct {
	AggValOID OID // type for aggVal
}

type FunctionMeta struct {
	ProNamespace    OID   // Namespace OID
	ProOwner        int32 // Owner's user ID
	ProLang         OID   // Language OID (e.g., SQL, PL/pgSQL, etc.)
	ProCost         int32 // Estimated execution cost
	ProRows         int32 // Estimated number of rows returned by the function
	ProVariadic     int32 // Is this a variadic function?
	ProSupport      int32
	ProKind         byte  // 'f' for normal function, 'p' for procedure, 'a' for aggregate, 'w' for window function
	ProSecDef       bool  // Security definer?
	ProLeakproof    bool  // Is the function leak-proof?
	ProIsStrict     bool  // Does the function return NULL when any input is NULL?
	ProVolatile     byte  // 'i' for immutable, 's' for stable, 'v' for volatile
	ProParallel     byte  // 'u' for unsafe, 's' for safe, 'r' for restricted
	ProNArgs        int16 // Number of arguments
	ProNArgDefaults int16 // Number of default arguments (from the right)

	ProAllArgTypes []OID    // All argument types (including variadic)
	ProArgModes    []byte   // Argument modes ('i' for IN, 'o' for OUT, 'b' for INOUT, 'v' for VARIADIC, 't' for TABLE)
	ProArgNames    []string // Argument names
	ProArgDefaults []string // Default argument values (as strings)
	ProTrfTypes    []OID    // Transform types

	ProSrc     string // Source code of the function (Depends on language) For SQL functions, this is the SQL definition; for PL/pgSQL, this is the function body.
	ProBin     []byte // Binary representation (if applicable)
	ProSQLBody string
	ProConfig  string //
	ProACL     string // Access control list (as a string representation)
}

// FunctionRegistry manages registered user-defined functions
type FunctionRegistry struct {
	mu              sync.RWMutex
	functions       map[OID]IFunction      // Функции по OID
	functionsByName map[string][]IFunction // Функции по имени (для перегрузки)
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions:       make(map[OID]IFunction),
		functionsByName: make(map[string][]IFunction),
	}
}

// Register adds a function to the registry
func (r *FunctionRegistry) Register(fn IFunction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	funcDef := fn.GetFunction()

	// Проверяем, не зарегистрирован ли уже этот OID
	if _, exists := r.functions[funcDef.OID]; exists {
		return fmt.Errorf("function %s already registered with OID %d", funcDef.ProName, funcDef.OID)
	}

	// Регистрируем по OID
	r.functions[funcDef.OID] = fn

	// Добавляем в индекс по имени для поддержки перегрузки
	r.functionsByName[funcDef.ProName] = append(r.functionsByName[funcDef.ProName], fn)

	return nil
}

// Get retrieves a function by OID
func (r *FunctionRegistry) Get(oid OID) (IFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, exists := r.functions[oid]
	return fn, exists
}

// GetByName retrieves first function by name (без учета перегрузки)
// Для поиска с учетом типов аргументов используйте GetByNameAndArgTypes
func (r *FunctionRegistry) GetByName(name string) (IFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	funcs, exists := r.functionsByName[name]
	if !exists || len(funcs) == 0 {
		return nil, false
	}

	// Возвращаем первую найденную функцию
	return funcs[0], true
}

// GetByNameAndArgTypes retrieves a function by name and argument types (поддержка перегрузки)
func (r *FunctionRegistry) GetByNameAndArgTypes(name string, argTypes []OID) (IFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	funcs, exists := r.functionsByName[name]
	if !exists || len(funcs) == 0 {
		return nil, false
	}

	// Ищем точное соответствие типов аргументов
	for _, fn := range funcs {
		funcDef := fn.GetFunction()
		if len(funcDef.ProArgTypes) == len(argTypes) {
			match := true
			for i, argType := range argTypes {
				if funcDef.ProArgTypes[i] != argType {
					match = false
					break
				}
			}
			if match {
				return fn, true
			}
		}
	}

	// Если точное совпадение не найдено, возвращаем первую функцию (для обратной совместимости)
	return funcs[0], true
}

// GetAllByName retrieves all function overloads by name
func (r *FunctionRegistry) GetAllByName(name string) []IFunction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	funcs, exists := r.functionsByName[name]
	if !exists {
		return nil
	}

	// Возвращаем копию среза
	result := make([]IFunction, len(funcs))
	copy(result, funcs)
	return result
}

// Unregister removes a function from the registry
func (r *FunctionRegistry) Unregister(oid OID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Получаем функцию перед удалением
	fn, exists := r.functions[oid]
	if !exists {
		return
	}

	funcDef := fn.GetFunction()

	// Удаляем из основного хранилища
	delete(r.functions, oid)

	// Удаляем из индекса по имени
	if funcs, exists := r.functionsByName[funcDef.ProName]; exists {
		// Ищем и удаляем эту функцию из списка
		for i, f := range funcs {
			if f.GetFunction().OID == oid {
				// Удаляем элемент из слайса
				r.functionsByName[funcDef.ProName] = append(funcs[:i], funcs[i+1:]...)

				// Если это была последняя функция с таким именем, удаляем запись
				if len(r.functionsByName[funcDef.ProName]) == 0 {
					delete(r.functionsByName, funcDef.ProName)
				}
				break
			}
		}
	}
}
