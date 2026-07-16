package jscontext

import (
	"context"

	"github.com/xmx/jsvm"
)

func New() jsvm.ModuleExporter { return new(contextModule) }

type contextModule struct{}

func (m *contextModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		"background":   context.Background,
		"withCancel":   context.WithCancel,
		"withTimeout":  context.WithTimeout,
		"withValue":    context.WithValue,
		"withDeadline": context.WithDeadline,
	}

	return jsvm.ModuleExports{
		Name:    "context",
		Default: defaults,
	}
}
