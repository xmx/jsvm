package jsvm

import "github.com/dop251/goja"

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
