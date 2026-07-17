package jsruntime

import (
	"runtime"

	"github.com/xmx/jsvm"
)

// runtimeModule 是 runtime 模块的实现，将 Go 运行时信息导出到 JS。
type runtimeModule struct{}

// New 创建 runtime 模块实例。
func New() jsvm.ModuleExporter {
	return new(runtimeModule)
}

// ModuleExports 注册 runtime 模块，导出：
//   - GOOS：当前操作系统（如 "linux"、"windows"、"darwin"）
//   - GOARCH：当前 CPU 架构（如 "amd64"、"arm64"）
//   - Compiler：编译器名称（如 "gc"）
//   - readMemStats：读取内存使用统计
//   - numCgoCall：cgo 调用次数
//   - numCPU：可用 CPU 核心数
//   - numGoroutine：当前存活 goroutine 数量
//   - version：Go 版本号
func (m *runtimeModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		"readMemStats": m.readMemStats,       // 读取内存统计（会触发 GC）
		"GOOS":         runtime.GOOS,         // 操作系统
		"GOARCH":       runtime.GOARCH,       // CPU 架构
		"Compiler":     runtime.Compiler,     // 编译器
		"numCgoCall":   runtime.NumCgoCall,   // cgo 调用次数
		"numCPU":       runtime.NumCPU,       // CPU 核心数
		"numGoroutine": runtime.NumGoroutine, // 当前 goroutine 数
		"version":      runtime.Version,      // Go 版本
	}

	return jsvm.ModuleExports{
		Name:    "runtime",
		Default: defaults,
	}
}

// readMemStats 读取当前 Go 运行时的内存使用统计信息。
// 调用前会先触发一次 GC，返回完整的 MemStats 结构。
func (m *runtimeModule) readMemStats() *runtime.MemStats {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)

	return stats
}
