package jsvm

import (
	"maps"
)

// ModuleExporter is what a module needs to return
type ModuleExporter interface {
	ModuleExports(vm *VM) ModuleExports
}

// ModuleExports is representation of ESM exports of a module
type ModuleExports struct {
	Name string

	// Default is what will be the `default` export of a module
	Default any

	// Named is the named exports of a module
	Named map[string]any
}

func (exp ModuleExports) toESM() any {
	if exp.Named == nil {
		return exp.Default
	}
	if exp.Default == nil {
		return exp.Named
	}

	result := make(map[string]any, len(exp.Named)+2)
	maps.Copy(result, exp.Named)
	result["default"] = exp.Default
	// This is to interop with any code that is transpiled by Babel or any similar tool.
	result["__esModule"] = true

	return result
}
