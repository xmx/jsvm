package jsruntime

import (
	"runtime"

	"github.com/xmx/jsvm"
)

type runtimeModule struct{}

func New() jsvm.ModuleExporter {
	return new(runtimeModule)
}

func (m *runtimeModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		"readMemStats": m.readMemStats,
		"GOOS":         runtime.GOOS,
		"GOARCH":       runtime.GOARCH,
		"Compiler":     runtime.Compiler,
		"numCgoCall":   runtime.NumCgoCall,
		"numCPU":       runtime.NumCPU,
		"numGoroutine": runtime.NumGoroutine,
		"version":      runtime.Version,
	}

	return jsvm.ModuleExports{
		Name:    "runtime",
		Default: defaults,
	}
}

func (m *runtimeModule) readMemStats() *runtime.MemStats {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)

	return stats
}
