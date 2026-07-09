package jsvm

import "github.com/dop251/goja"

// LoadModuleOptions configures module loading. Reserved for future use.
type LoadModuleOptions struct {
	Unsafe bool
}

// ModuleLoader loads modules into the VM's require registry.
type ModuleLoader interface {
	// LoadModule loads a module and returns its package name and exported values.
	LoadModule(vm *VM, opts LoadModuleOptions) (pkg string, vals map[string]any, err error)
}

type Module interface {
	// Name returns the module name.
	Name() string

	// Load loads the module.
	Load(vm *VM, exports *goja.Object) error
}

func SetExports(exports *goja.Object, vals map[string]any) error {
	for k, v := range vals {
		if err := exports.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}
