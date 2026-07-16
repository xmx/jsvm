package jshttp

import (
	"net/http"
	"sync"
)

type serialHandler struct {
	m sync.Mutex
	h http.Handler
}

func (s *serialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.Lock()
	defer s.m.Unlock()

	s.h.ServeHTTP(w, r)
}

func (m *httpModule) listenAndServe(addr string, handler http.Handler) error {
	if handler == nil {
		handler = http.NotFoundHandler()
	}

	safe := &serialHandler{h: handler}
	srv := &http.Server{
		Addr:    addr,
		Handler: safe,
	}
	cln, ok := m.vm.AddCleaner(srv)
	if !ok {
		return http.ErrServerClosed
	}
	err := srv.ListenAndServe()
	_ = cln.Close()

	return err
}
