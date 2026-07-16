package jshttp

import (
	"io"
	"net/http"
	"strings"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

type httpModule struct {
	vm *jsvm.VM
}

func New() jsvm.ModuleExporter {
	return new(httpModule)
}

func (m *httpModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm
	defaults := map[string]any{
		"fetch":          m.fetch,
		"newServeMux":    http.NewServeMux,
		"listenAndServe": m.listenAndServe,
	}

	return jsvm.ModuleExports{
		Name:    "net/http",
		Default: defaults,
	}
}

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
			req.Host = v
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

type requestInfo struct {
	Method string
	Header map[string]string
	Body   io.Reader
}
