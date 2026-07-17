package jsvm

import (
	"maps"
)

// ModuleExporter 是模块需要实现的接口，用于向 VM 导出模块内容。
type ModuleExporter interface {
	// ModuleExports 返回该模块的导出定义。
	// vm 是该模块即将绑定的 VM 实例，可用于访问 VM 的上下文、日志等。
	ModuleExports(vm *VM) ModuleExports
}

// ModuleExports 代表一个模块的 ESM 导出结构。
type ModuleExports struct {
	// Name 是模块的注册名称，脚本 import 时使用此名称，如 "net/http"。
	Name string

	// Default 是模块的默认导出（对应 JS 的 export default），通常为对象或函数。
	Default any

	// Named 是模块的命名导出集合（对应 JS 的 export const/let），键为导出名。
	Named map[string]any
}

// toESM 将 ModuleExports 转换为供 require 返回的 ESM 兼容对象。
// 仅有 Default 时直接返回 Default；
// 仅有 Named 时直接返回 Named；
// 两者都有时构造一个包含所有命名导出、default 和 __esModule 标记的对象。
func (exp ModuleExports) toESM() any {
	if exp.Named == nil {
		return exp.Default
	}
	if exp.Default == nil {
		return exp.Named
	}

	// 同时存在默认导出和命名导出，构造完整的 ESM 对象
	result := make(map[string]any, len(exp.Named)+2)
	maps.Copy(result, exp.Named)
	result["default"] = exp.Default
	// __esModule 标记用于与 Babel 等工具转译后的代码兼容
	result["__esModule"] = true

	return result
}
