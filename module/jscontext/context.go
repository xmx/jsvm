package jscontext

import (
	"context"

	"github.com/xmx/jsvm"
)

// New 创建 context 模块实例。
func New() jsvm.ModuleExporter { return new(contextModule) }

// contextModule 是 context 模块的实现，将 Go 标准库 context 包的主要函数导出到 JS。
type contextModule struct{}

// ModuleExports 注册 context 模块，导出 background/withCancel/withTimeout/withValue/withDeadline。
func (m *contextModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		"background":   context.Background,   // 创建永不过期的根上下文
		"withCancel":   context.WithCancel,   // 创建可手动取消的子上下文
		"withTimeout":  context.WithTimeout,  // 创建带超时的子上下文
		"withValue":    context.WithValue,    // 创建携带键值对的子上下文
		"withDeadline": context.WithDeadline, // 创建带截止时间的子上下文
	}

	return jsvm.ModuleExports{
		Name:    "context",
		Default: defaults,
	}
}
