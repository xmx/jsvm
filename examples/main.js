import console from 'console'
import http from 'net/http'

async function foo() {
    return "111"
}

const mux = http.newServeMux()
mux.handleFunc('/ping', async (w, r) => {
    console.log('收到请求 ' + r.remoteAddr)
    const s = await foo()
    w.write(s)
})

http.listenAndServe(':8088', mux)