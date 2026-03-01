package baseplugin

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	basefuncs "petacore/plugins/base/funcs"
	psdk "petacore/sdk"

	"go.uber.org/zap"
)

type BasePlugin struct {
	config map[string]interface{}
}

func (p *BasePlugin) Name() string {
	return "base"
}

func (p *BasePlugin) Version() psdk.Version {
	return psdk.NewVersion(0, 0, 1)
}

// Init initializes the plugin
func (p *BasePlugin) Init(config map[string]interface{}) error {
	p.config = config
	fmt.Println("BasePlugin initialized with config:", config)
	return nil
}

// Execute executes the plugin
func (p *BasePlugin) Execute(ctx context.Context, args ...interface{}) (interface{}, error) {
	logger.Info("BasePlugin executed with args:", zap.Any("args:", args))
	return nil, nil
}

// Shutdown shuts down the plugin
func (p *BasePlugin) Shutdown() error {
	logger.Info("BasePlugin shutdown")
	return nil
}

// RegisterFunctions registers the plugin's functions
func RegisterFunctions(registry *psdk.FunctionRegistry) error {
	funcs := []psdk.IFunction{
		&UpperFunction{},
		// Агрегатные функции
		&CountFunction{},
		&SumFunction{},
		&AvgFunction{},
		// MAX с разными типами (перегрузка)
		&MaxFunction{},     // float8
		&MaxFunctionInt{},  // int4
		&MaxFunctionText{}, // text
		// MIN с разными типами (перегрузка)
		&MinFunction{},     // float8
		&MinFunctionInt{},  // int4
		&MinFunctionText{}, // text
		// Обычные функции
		&basefuncs.CurrentDatabaseFunction{},
		&basefuncs.CurrentCatalogFunction{},
		&basefuncs.CurrentSchemaFunction{},
		&basefuncs.CurrentSchemasFunction{},
		&basefuncs.CurrentUserFunction{},
		&basefuncs.CurrentRoleFunction{},
		&basefuncs.SessionUserFunction{},
		&basefuncs.UserFunction{},
		&basefuncs.NowFunction{},
		&basefuncs.QuoteIdentFunction{},
		&basefuncs.SubstringFunction{},
		&basefuncs.PgTableIsVisibleFunction{},
		&basefuncs.LengthFunction{},
		&basefuncs.VersionFunction{},
		&basefuncs.PgBackendPidFunction{},
		&basefuncs.PgPostmasterStartTimeFunction{},
		// Функции для работы с массивами
		&basefuncs.ArrayToStringFunction{},
		&basefuncs.ArrayToStringIntFunction{},
		&basefuncs.ArrayLengthFunction{},
		&basefuncs.ArrayLengthIntFunction{},
	}

	logger.Debugf("Registering functions for plugin %s: %v\n", Plugin.Name(), funcs)

	for _, fn := range funcs {
		if err := registry.Register(fn); err != nil {
			return fmt.Errorf("failed to register function %s: %w", fn.GetFunction().ProName, err)
		}
	}
	return nil
}

// Plugin is the exported plugin instance
var Plugin psdk.PetaPlugin = &BasePlugin{}
