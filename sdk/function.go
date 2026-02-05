package psdk

import (
	"context"
	"fmt"
	"sync"
)

type IFunction interface {
	Execute(ctx context.Context, args ...any) (any, error)
	GetFunction() *Function
}

// Function represents a user-defined function
type Function struct {
	OID         OID    // Unique identifier for the function
	ProName     string // Function name
	ProArgTypes []OID  // Argument types (as OIDs)
	ProRetType  OID    // Return type (as OID)

	Meta FunctionMeta // Additional metadata about the function
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
	mu        sync.RWMutex
	functions map[OID]IFunction
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[OID]IFunction),
	}
}

// Register adds a function to the registry
func (r *FunctionRegistry) Register(fn IFunction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.functions[fn.GetFunction().OID]; exists {
		return fmt.Errorf("function %s already registered with OID %d", fn.GetFunction().ProName, fn.GetFunction().OID)
	}
	r.functions[fn.GetFunction().OID] = fn
	return nil
}

// Get retrieves a function by OID
func (r *FunctionRegistry) Get(oid OID) (IFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, exists := r.functions[oid]
	return fn, exists
}

func (r *FunctionRegistry) GetByName(name string) (IFunction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, fn := range r.functions {
		if fn.GetFunction().ProName == name {
			return fn, true
		}
	}
	return nil, false
}

// Unregister removes a function from the registry
func (r *FunctionRegistry) Unregister(oid OID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.functions, oid)
}
