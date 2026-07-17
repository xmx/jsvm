package jshttp

import (
	"io"
	"net/http"
	"strings"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

// httpModule 是 net/http 模块的实现，提供 HTTP 客户端和服务端功能。
type httpModule struct {
	vm *jsvm.VM
}

// New 创建 net/http 模块实例。
func New() jsvm.ModuleExporter {
	return new(httpModule)
}

// ModuleExports 注册 net/http 模块，导出 fetch、newServeMux、listenAndServe。
func (m *httpModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm
	defaults := map[string]any{
		"fetch":          m.fetch,          // 发起 HTTP 请求
		"newServeMux":    http.NewServeMux, // 创建 HTTP 请求路由器
		"listenAndServe": m.listenAndServe, // 启动 HTTP 服务器
	}

	return jsvm.ModuleExports{
		Name:    "net/http",
		Default: defaults,
	}
}

// fetch 发起 HTTP 请求并返回封装的 Response 对象。
// 参数：fetch(url, {method?, header?, body?})
func (m *httpModule) fetch(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	strURL := call.Argument(0).String()
	info := m.parseRequestInfo(call.Argument(1), rt)

	ctx := m.vm.Context()
	req, err := http.NewRequestWithContext(ctx, info.Method, strURL, info.Body)
	if err != nil {
		m.vm.Throw(err)
	}
	for k, v := range info.Header {
		req.Header.Set(k, v)
		if k == "host" {
			req.Host = v // Host 需单独设置（Go 规定）
		}
	}

	cli := http.DefaultClient
	res, err1 := cli.Do(req)
	if err != nil {
		m.vm.Throw(err1)
	}
	resp := &response{Response: res}

	return rt.ToValue(resp)
}

// parseRequestInfo 从 JS 对象中解析出请求配置（方法、请求头、请求体）。
// 若 val 为 nullish 则返回默认配置（GET 方法，空请求头）。
func (m *httpModule) parseRequestInfo(val goja.Value, rt *goja.Runtime) requestInfo {
	info := requestInfo{Header: make(map[string]string)}
	if jsvm.IsNullish(val) {
		return info
	}

	obj := val.ToObject(rt)
	{
		v := obj.Get("method")
		if !jsvm.IsNullish(v) {
			info.Method = v.String()
		}
	}
	{
		v := obj.Get("header")
		if !jsvm.IsNullish(v) {
			h := v.ToObject(rt)
			for _, key := range h.Keys() {
				str := h.Get(key).String()
				// 请求头名统一小写存储，避免大小写不一致问题
				info.Header[strings.ToLower(key)] = str
			}
		}
	}
	{
		v := obj.Get("body")
		if !jsvm.IsNullish(v) {
			info.Body = strings.NewReader(v.String())
		}
	}

	return info
}

// requestInfo 是发起 HTTP 请求所需的内部配置结构。
type requestInfo struct {
	Method string            // HTTP 方法（GET/POST/PUT/DELETE 等，不传默认为空）
	Header map[string]string // 自定义请求头（键为小写）
	Body   io.Reader         // 请求体内容
}
