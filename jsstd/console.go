package jsstd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/xmx/jsvm"
)

type consoleModule struct {
	vm  *jsvm.VM
	log *slog.Logger
}

func NewConsole() jsvm.ModuleLoader {
	return &consoleModule{}
}

func (m *consoleModule) LoadModule(vm *jsvm.VM, opts jsvm.LoadModuleOptions) (string, map[string]any, error) {
	m.vm = vm
	m.log = vm.Logger()
	vals := map[string]any{
		"log":   m.logFunc(slog.LevelInfo),
		"debug": m.logFunc(slog.LevelDebug),
		"warn":  m.logFunc(slog.LevelWarn),
		"error": m.logFunc(slog.LevelError),
	}

	return "console", vals, nil
}

func (m *consoleModule) logFunc(level slog.Level) func(...any) {
	return func(args ...any) {
		parts := make([]string, 0, len(args))
		for _, arg := range args {
			parts = append(parts, fmt.Sprint(arg))
		}
		msg := strings.Join(parts, " ")

		switch level {
		case slog.LevelDebug:
			m.log.Debug(msg)
		case slog.LevelWarn:
			m.log.Warn(msg)
		case slog.LevelError:
			m.log.Error(msg)
		default:
			m.log.Info(msg)
		}
	}
}
