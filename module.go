package jsvm

// LoadModuleOptions configures module loading. Reserved for future use.
type LoadModuleOptions struct {
	Unsafe bool
}

// ModuleLoader loads modules into the VM's require registry.
type ModuleLoader interface {
	// LoadModule loads a module and returns its package name and exported values.
	LoadModule(vm *VM, opts LoadModuleOptions) (pkg string, vals map[string]any, err error)
}
