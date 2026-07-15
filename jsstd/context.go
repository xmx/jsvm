package jsstd

import (
	"context"

	"github.com/xmx/jsvm"
)

func NewContext() jsvm.ModuleExporter { return new(contextModule) }

type contextModule struct{}

func (m *contextModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	vals := map[string]any{
		"background":   context.Background,
		"withCancel":   context.WithCancel,
		"withTimeout":  context.WithTimeout,
		"withValue":    context.WithValue,
		"withDeadline": context.WithDeadline,
	}

	return jsvm.ModuleExports{
		Name:    "context",
		Default: vals,
	}
}
