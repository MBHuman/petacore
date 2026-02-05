package baseplugin

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	psdk "petacore/sdk"
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
	fmt.Println("BasePlugin executed with args:", args)
	return "BasePlugin result: " + fmt.Sprintf("%v", args), nil
}

// Shutdown shuts down the plugin
func (p *BasePlugin) Shutdown() error {
	fmt.Println("BasePlugin shutdown")
	return nil
}

// RegisterFunctions registers the plugin's functions
func RegisterFunctions(registry *psdk.FunctionRegistry) error {
	funcs := []psdk.IFunction{
		&UpperFunction{},
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
