package jsvm

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/parser"
	"github.com/evanw/esbuild/pkg/api"
)

// Transform transpiles the input source string and strip types from it.
// this is done using esbuild
//
// https://github.com/grafana/k6/blob/c0bc819af0fdb3595fbafc6897cd8abd32da9d84/internal/js/compiler/enhanced.go#L9
func Transform(name, code string) api.TransformResult {
	opts := api.TransformOptions{
		LogLevel:      api.LogLevelSilent,
		Sourcemap:     api.SourceMapInline,
		Target:        api.ES2017,
		Platform:      api.PlatformNeutral,
		Format:        api.FormatCommonJS,
		Charset:       api.CharsetUTF8,
		LegalComments: api.LegalCommentsNone,
		Sourcefile:    name,
		Loader:        api.LoaderJS,
	}

	return api.Transform(code, opts)
}

// Compile transpiles and compiles JavaScript code into a goja program.
func Compile(name, code string) (*goja.Program, error) {
	res := Transform(name, code)
	if len(res.Errors) == 0 {
		cjs := string(res.Code)
		return goja.Compile(name, cjs, false)
	}

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
