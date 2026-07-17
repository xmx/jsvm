package jsconsole

import (
	"bytes"
	"log/slog"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

// consoleModule 是 console 模块的实现，将 Go 的 slog 日志输出到 JS 脚本中。
type consoleModule struct {
	vm *jsvm.VM
}

// New 创建 console 模块实例。
func New() jsvm.ModuleExporter {
	return &consoleModule{}
}

// print 返回一个 JS 可调用函数，将传入参数格式化后按指定日志级别输出。
// 若当前日志级别低于 lvl，则静默返回 undefined。
func (m *consoleModule) print(lvl slog.Level) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		ctx := m.vm.Context()
		log := m.vm.Logger()
		if !log.Enabled(ctx, lvl) {
			return goja.Undefined() // 日志级别未开启，直接返回
		}

		buf := m.format(call)
		msg := buf.String()
		log.Log(ctx, lvl, msg)

		return goja.Undefined()
	}
}

// ModuleExports 注册 console 模块，将 log/info/debug/warn/error 函数导出。
func (m *consoleModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm
	defaults := map[string]any{
		"log":   m.print(slog.LevelInfo),
		"info":  m.print(slog.LevelInfo),
		"debug": m.print(slog.LevelDebug),
		"warn":  m.print(slog.LevelWarn),
		"error": m.print(slog.LevelError),
	}

	return jsvm.ModuleExports{
		Name:    "console",
		Default: defaults,
	}
}

// format 将 JS 函数调用的所有参数格式化为单个字符串，参数间以空格分隔。
func (m *consoleModule) format(call goja.FunctionCall) *bytes.Buffer {
	dst := new(bytes.Buffer)
	num := len(call.Arguments)
	for i, arg := range call.Arguments {
		m.parseAny(dst, arg)
		if i < num-1 {
			dst.WriteByte(' ') // 参数之间插入空格分隔
		}
	}

	return dst
}

// parseAny 将单个 goja 值转换为字符串写入缓冲区，使用 StringValue 进行类型感知转换。
func (m *consoleModule) parseAny(dst *bytes.Buffer, val goja.Value) {
	msg := jsvm.StringValue(val)
	dst.WriteString(msg)
}
