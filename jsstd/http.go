package jsstd

import (
	"context"
	"net/http"

	"github.com/xmx/jsvm"
)

type httpModule struct {
	vm *jsvm.VM
}

func NewHTTP() jsvm.ModuleExporter {
	return &httpModule{}
}

func (m *httpModule) ModuleExports(vm *jsvm.VM) jsvm.ModuleExports {
	m.vm = vm
	vals := map[string]any{
		"newServeMux":    http.NewServeMux,
		"listenAndServe": m.listenAndServe,
	}

	return jsvm.ModuleExports{
		Name:    "net/http",
		Default: vals,
	}
}

func (m *httpModule) listenAndServe(addr string, h http.Handler) error {
	if h == nil {
		h = http.NotFoundHandler()
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	cln, succ := m.vm.AddCleaner(srv)
	if !succ {
		return context.Canceled
	}

	err := srv.ListenAndServe()
	_ = cln.Close()

	return err
}
