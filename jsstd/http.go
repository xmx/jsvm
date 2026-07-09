package jsstd

import (
	"net/http"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

type httpModule struct {
	vm *jsvm.VM
}

func NewHTTP() jsvm.Module {
	return &httpModule{}
}

func (m *httpModule) Name() string {
	return "net/http"
}

func (m *httpModule) Load(vm *jsvm.VM, exports *goja.Object) error {
	m.vm = vm
	vals := map[string]any{
		"newServeMux":    http.NewServeMux,
		"listenAndServe": m.listenAndServe,
	}

	return jsvm.SetExports(exports, vals)
}

func (m *httpModule) listenAndServe(addr string, h http.Handler) error {
	if h == nil {
		h = http.NotFoundHandler()
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	cln, _ := m.vm.AddCleaner(srv)
	err := srv.ListenAndServe()
	_ = cln.Close()

	return err
}
