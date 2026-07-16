import console from 'console'
import runtime from 'runtime'
import http from 'net/http'
import httputil from 'net/http/httputil'

const pxy = new httputil.ReverseProxy('https://mirrors.zju.edu.cn/')
const mux = http.newServeMux()
mux.handleFunc('/', (w, r) => {
    pxy.serveHTTP(w, r)
})

mux.handleFunc('/stats', (w, r) => {
    const stats = runtime.readMemStats()
    const str = JSON.stringify(stats)
    w.header().set('Content-Type', 'application/json')
    w.write(str)
})

http.listenAndServe(':9999', mux)
