package jsstd

import (
	"io"

	"github.com/xmx/jsvm"
)

type ioModule struct{}

func NewIO() jsvm.ModuleExporter {
	return &httpModule{}
}

func (m *ioModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	vals := map[string]any{
		"copy":        io.Copy,
		"discard":     io.Discard,
		"EOF":         io.EOF,
		"limitReader": io.LimitReader,
	}

	return jsvm.ModuleExports{
		Name:    "io",
		Default: vals,
	}
}
