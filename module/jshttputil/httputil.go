package jshttputil

import (
	"net/http/httputil"
	"net/url"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

// httputilModule 是 net/http/httputil 模块的实现，提供 HTTP 反向代理等工具。
type httputilModule struct {
	vm *jsvm.VM
}

// New 创建 net/http/httputil 模块实例。
func New() jsvm.ModuleExporter {
	return new(httputilModule)
}

// ModuleExports 注册 net/http/httputil 模块，导出 ReverseProxy 构造函数。
func (m *httputilModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm

	defaults := map[string]any{
		"ReverseProxy": m.newReverseProxy, // 创建反向代理
	}

	return jsvm.ModuleExports{
		Name:    "net/http/httputil",
		Default: defaults,
	}
}

// newReverseProxy 创建指向指定目标的反向代理，实现 http.Handler 接口。
// 参数：target 可以是 URL 字符串或 url.URL 对象。
func (m *httputilModule) newReverseProxy(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object {
	pxy := new(httputil.ReverseProxy)
	arg0 := call.Argument(0)
	if !jsvm.IsNullish(arg0) {
		exp := arg0.Export()
		switch uv := exp.(type) {
		case *url.URL:
			// 已解析的 URL 对象，直接创建代理
			pxy = httputil.NewSingleHostReverseProxy(uv)
		default:
			// 字符串 URL，先解析再创建代理
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
