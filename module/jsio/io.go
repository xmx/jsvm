package jsio

import (
	"io"

	"github.com/xmx/jsvm"
)

// ioModule 是 io 模块的实现，将 Go 标准库 io 包的基本读写工具导出到 JS。
type ioModule struct{}

// New 创建 io 模块实例。
func New() jsvm.ModuleExporter {
	return &ioModule{}
}

// ModuleExports 注册 io 模块，导出：
//   - copy：将 Reader 数据写入 Writer
//   - discard：丢弃写入数据的 Writer
//   - EOF：流结束哨兵错误
//   - limitReader：创建最多读取 n 字节的有限 Reader
//   - readAll：读取 Reader 全部内容
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
