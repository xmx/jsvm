# net/http/httputil

## 模块内容

```typescript
interface httputil {
    /**
     * ReverseProxy 构造函数，创建 HTTP 反向代理
     * @param target 目标服务器地址，可以是 URL 字符串或 url.URL 对象
     */
    ReverseProxy: {
        new(target: string | URL): ReverseProxy;
    };
}
```

### ReverseProxy

`ReverseProxy` 是 `http.Handler` 的实现，将收到的请求转发到指定的目标地址，并将目标响应返回给客户端。

```typescript
interface ReverseProxy extends Handler {
    // 继承 http.Handler 接口，可直接传给 listenAndServe
}
```

## 使用样例

```javascript
import console from 'console'
import http from 'net/http'
import httputil from 'net/http/httputil'

// 创建反向代理，将所有请求转发到目标服务器
const proxy = new httputil.ReverseProxy('http://localhost:3000')

console.log('Proxy server starting on :8080')
http.listenAndServe(':8080', proxy)
```
