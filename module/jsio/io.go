package jsio

import (
	"io"

	"github.com/xmx/jsvm"
)

type ioModule struct{}

func New() jsvm.ModuleExporter {
	return &ioModule{}
}

func (m *ioModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	vals := map[string]any{
		"copy":        io.Copy,
		"discard":     io.Discard,
		"EOF":         io.EOF,
		"limitReader": io.LimitReader,
		"readAll":     io.ReadAll,
	}

	return jsvm.ModuleExports{
		Name:    "io",
		Default: vals,
	}
}
