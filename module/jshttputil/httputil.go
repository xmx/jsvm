package jshttputil

import (
	"net/http/httputil"
	"net/url"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

type httputilModule struct {
	vm *jsvm.VM
}

func New() jsvm.ModuleExporter {
	return new(httputilModule)
}

func (m *httputilModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm

	defaults := map[string]any{
		"ReverseProxy": m.newReverseProxy,
	}

	return jsvm.ModuleExports{
		Name:    "net/http/httputil",
		Default: defaults,
	}
}

func (m *httputilModule) newReverseProxy(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object {
	pxy := new(httputil.ReverseProxy)
	arg0 := call.Argument(0)
	if !jsvm.IsNullish(arg0) {
		exp := arg0.Export()
		switch uv := exp.(type) {
		case *url.URL:
			pxy = httputil.NewSingleHostReverseProxy(uv)
		default:
			strURL := arg0.String()
			pu, err := url.Parse(strURL)
			if err != nil {
				m.vm.Throw(err)
			}
			pxy = httputil.NewSingleHostReverseProxy(pu)
		}
	}

	return rt.ToValue(pxy).ToObject(rt)
}
