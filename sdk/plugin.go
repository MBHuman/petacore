package psdk

import (
	"context"
	"fmt"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func NewVersion(major, minor, patch int) Version {
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func NewVersionFromString(versionStr string) (Version, error) {
	var v Version
	_, err := fmt.Sscanf(versionStr, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	if err != nil {
		return Version{}, fmt.Errorf("invalid version format: %s", versionStr)
	}
	return v, nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Cmp(other Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// PetaPlugin defines the interface for plugins in PetaCore
type PetaPlugin interface {
	// Name returns the unique name of the plugin
	Name() string
	// Init initializes the plugin with configuration
	Init(config map[string]interface{}) error
	// Version returns the plugin version
	Version() Version
	// Execute runs the plugin's main logic
	Execute(ctx context.Context, args ...interface{}) (interface{}, error)
	// Shutdown cleans up resources
	Shutdown() error
}

// PetaPluginRegistry manages registered plugins
type PetaPluginRegistry struct {
	plugins map[string]PetaPlugin
}

// NewPetaPluginRegistry creates a new plugin registry
func NewPetaPluginRegistry() *PetaPluginRegistry {
	return &PetaPluginRegistry{
		plugins: make(map[string]PetaPlugin),
	}
}

// Register adds a plugin to the registry
func (r *PetaPluginRegistry) Register(plugin PetaPlugin) error {
	if _, exists := r.plugins[plugin.Name()]; exists {
		return fmt.Errorf("plugin %s already registered", plugin.Name())
	}
	r.plugins[plugin.Name()] = plugin
	return nil
}

// Get retrieves a plugin by name
func (r *PetaPluginRegistry) Get(name string) (PetaPlugin, bool) {
	plugin, exists := r.plugins[name]
	return plugin, exists
}

// Unregister removes a plugin from the registry
func (r *PetaPluginRegistry) Unregister(name string) {
	delete(r.plugins, name)
}
