package jsstd

import (
	"bytes"
	"log/slog"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

type consoleModule struct {
	vm *jsvm.VM
}

func NewConsole() jsvm.ModuleExporter {
	return &consoleModule{}
}

func (m *consoleModule) print(lvl slog.Level) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		ctx := m.vm.Context()
		log := m.vm.Logger()
		if !log.Enabled(ctx, lvl) {
			return goja.Undefined()
		}

		buf := m.format(call)
		msg := buf.String()
		log.Log(ctx, lvl, msg)

		return goja.Undefined()
	}
}

func (m *consoleModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm
	vals := map[string]any{
		"log":   m.print(slog.LevelInfo),
		"info":  m.print(slog.LevelInfo),
		"debug": m.print(slog.LevelDebug),
		"warn":  m.print(slog.LevelWarn),
		"error": m.print(slog.LevelError),
	}

	return jsvm.ModuleExports{
		Name:    "console",
		Default: vals,
	}
}

func (m *consoleModule) format(call goja.FunctionCall) *bytes.Buffer {
	dst := new(bytes.Buffer)
	num := len(call.Arguments)
	for i, arg := range call.Arguments {
		m.parseAny(dst, arg)
		if i < num-1 {
			dst.WriteByte(' ')
		}
	}

	return dst
}

func (m *consoleModule) parseAny(dst *bytes.Buffer, val goja.Value) {
	msg := jsvm.StringValue(val)
	dst.WriteString(msg)
}
