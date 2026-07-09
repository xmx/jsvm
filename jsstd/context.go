package jsstd

import (
	"context"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

func NewContext() jsvm.Module { return new(contextModule) }

type contextModule struct{}

func (m *contextModule) Name() string {
	return "context"
}

func (m *contextModule) Load(_ *jsvm.VM, exports *goja.Object) error {
	vals := map[string]any{
		"background":   context.Background,
		"withCancel":   context.WithCancel,
		"withTimeout":  context.WithTimeout,
		"withValue":    context.WithValue,
		"withDeadline": context.WithDeadline,
	}

	return jsvm.SetExports(exports, vals)
}
