package jsurl

import (
	"net/url"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

type urlModule struct{}

func New() jsvm.ModuleExporter {
	return &urlModule{}
}

func (m *urlModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		"parse":           url.Parse,
		"parseRequestURI": url.ParseRequestURI,
		"parseQuery":      url.ParseQuery,
		"queryEscape":     url.QueryEscape,
		"queryUnescape":   url.QueryUnescape,
		"user":            url.User,
		"userPassword":    url.UserPassword,
		"Values":          m.newValues,
	}

	return jsvm.ModuleExports{
		Name:    "net/url",
		Default: defaults,
	}
}

func (m *urlModule) newValues(_ goja.ConstructorCall, rt *goja.Runtime) *goja.Object {
	v := make(url.Values)
	val := rt.ToValue(v)

	return val.ToObject(rt)
}
