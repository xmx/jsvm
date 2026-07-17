package jshttp

import (
	"net/http"
	"sync"
)

// serialHandler 是同步化的 HTTP 处理器，对底层 handler 加全局互斥锁，
// 保证同时只有一个请求被执行（单线程处理）。适合避免 VM 并发问题。
type serialHandler struct {
	m sync.Mutex // 保护 h.ServeHTTP 的并发调用
	h http.Handler
}

// ServeHTTP 实现 http.Handler 接口，加锁后转发到底层处理器。
func (s *serialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.Lock()
	defer s.m.Unlock()

	s.h.ServeHTTP(w, r)
}

// listenAndServe 在指定地址启动 HTTP 服务器。
// 对 handler 加串化保护（防止 VM 并发调用），并将服务器注册到 VM 清理列表。
// VM 关闭时自动中止服务器。
func (m *httpModule) listenAndServe(addr string, handler http.Handler) error {
	if handler == nil {
		handler = http.NotFoundHandler()
	}

	// 包装为串行处理器，避免 goja 并发访问问题
	safe := &serialHandler{h: handler}
	srv := &http.Server{
		Addr:    addr,
		Handler: safe,
	}
	cln, ok := m.vm.AddCleaner(srv)
	if !ok {
		return http.ErrServerClosed // VM 已关闭，拒绝启动
	}
	err := srv.ListenAndServe()
	_ = cln.Close() // 服务器退出后触发资源清理

	return err
}
