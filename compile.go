package jsvm

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/parser"
	"github.com/evanw/esbuild/pkg/api"
)

// Transform 使用 esbuild 将输入源码转译为 CommonJS。
//
// 参考：https://github.com/grafana/k6/blob/c0bc819af0fdb3595fbafc6897cd8abd32da9d84/internal/js/compiler/enhanced.go#L9
func Transform(name, code string) api.TransformResult {
	opts := api.TransformOptions{
		LogLevel:      api.LogLevelSilent,    // 静默日志，错误通过结果返回
		Sourcemap:     api.SourceMapInline,   // 内联 Sourcemap，便于调试
		Target:        api.ES2017,            // 降级到 ES2017，goja 支持的上限
		Platform:      api.PlatformNeutral,   // 平台中立（非浏览器也非 Node）
		Format:        api.FormatCommonJS,    // 输出 CommonJS 格式
		Charset:       api.CharsetUTF8,       // UTF-8 编码
		LegalComments: api.LegalCommentsNone, // 不保留法律注释
		Sourcefile:    name,                  // 源文件名，用于错误提示
		Loader:        api.LoaderJS,          // 按 JS 加载器处理
	}

	return api.Transform(code, opts)
}

// Compile 将 JavaScript 源码转译并编译为 goja 可执行的程序。
// 先通过 Transform 处理，再用 goja.Compile 编译。
// name 是源文件名（用于错误信息），code 是源码字符串。
// 返回编译后的程序；若转译出错则返回包含详细行/列号的 ErrorList。
func Compile(name, code string) (*goja.Program, error) {
	res := Transform(name, code)
	if len(res.Errors) == 0 {
		cjs := string(res.Code)
		return goja.Compile(name, cjs, false)
	}

	// 将 esbuild 错误转换为 goja 的 ErrorList 格式
	var err parser.ErrorList
	for _, te := range res.Errors {
		var pos file.Position
		if l := te.Location; l != nil {
			pos.Filename = l.File
			pos.Line = l.Line
			pos.Column = l.Column
		}

		err.Add(pos, te.Text)
	}

	return nil, err
}
