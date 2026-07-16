# net/http

## 模块内容

```typescript
interface http {
    /**
     * 发起 HTTP 请求，返回响应对象
     * @param url 请求地址
     * @param info 可选的请求配置（方法、请求头、请求体）
     */
    fetch(url: string, info?: RequestInfo): Response;

    /**
     * 创建一个新的 HTTP 请求多路复用器
     */
    newServeMux(): ServeMux;

    /**
     * 启动 HTTP 服务器，监听指定地址
     * @param addr 监听地址，如 ":8080"
     * @param handler 可选的请求处理器，不传则返回 404
     */
    listenAndServe(addr: string, handler?: Handler): Error;
}
```

### RequestInfo

```typescript
/**
 * HTTP 请求配置
 */
interface RequestInfo {
    /** 请求方法，如 GET、POST、PUT 等 */
    method: string;
    /** 自定义请求头 */
    header: { [key: string]: string };
    /** 请求体内容 */
    body: string;
}
```

### Response

```typescript
/**
 * HTTP 响应对象
 */
interface Response {
    /** 响应状态文本，如 "200 OK" */
    readonly status: string;
    /** 响应状态码，如 200、404 */
    readonly statusCode: number;

    /** 读取响应体为文本字符串 */
    text(): string;

    /** 将响应体解析为 JSON 对象 */
    json(): any;
}
```

### Handler

```typescript
/**
 * HTTP 请求处理器接口
 */
interface Handler {
    /** 处理 HTTP 请求 */
    serveHTTP(w: ResponseWriter, r: Request): void;
}
```

### ResponseWriter

```typescript
/**
 * HTTP 响应写入器，用于构建和发送 HTTP 响应
 */
interface ResponseWriter {
    /** 响应头，可通过赋值设置响应头字段 */
    readonly header: { [key: string]: string };

    /** 向响应体写入数据 */
    write(data: Uint8Array | string): number;

    /** 设置响应状态码，必须在写入响应体之前调用 */
    writeHeader(statusCode: number): void;
}
```

### Request

```typescript
/**
 * HTTP 请求对象，包含客户端请求的全部信息
 */
interface Request {
    /** 请求方法，如 GET、POST */
    readonly method: string;
    /** 请求的 URL */
    readonly url: URL;
    /** 请求头 */
    readonly header: { [key: string]: string };
    /** 请求的 Host 头 */
    readonly host: string;
    /** 客户端远程地址 */
    readonly remoteAddr: string;

    /** 获取表单字段的值（支持 URL 参数和 POST 表单） */
    formValue(key: string): string;
}
```

### ServeMux

`ServeMux` 是 HTTP 请求多路复用器，实现了 `Handler` 接口。

```typescript
interface ServeMux extends Handler {
    /**
     * 注册路由处理器（需实现 Handler 接口）
     * @param pattern 路由模式，如 "/api/"
     * @param handler 请求处理器
     */
    handle(pattern: string, handler: Handler): void;

    /**
     * 注册路由处理函数
     * @param pattern 路由模式，如 "/api/"
     * @param handler 处理函数
     */
    handleFunc(pattern: string, handler: (w: ResponseWriter, r: Request) => void): void;
}
```

## 使用样例

### 发起 HTTP 请求

```javascript
import console from 'console'
import http from 'net/http'
import url from 'net/url'

const params = new url.Values()
params.set('host', '127.0.0.1')

const res = http.fetch('https://httpbin.io/get?' + params.encode(), {
    method: 'GET',
    header: {
        Accept: 'application/json',
        'User-Agent': 'Mozilla/5.0',
    },
});

const data = res.text()
console.log(data);
```

### 启动 HTTP 服务

```javascript
import console from 'console'
import http from 'net/http'

const mux = http.newServeMux()

mux.handleFunc('/hello', (w, r) => {
    w.write('hello world')
})

mux.handleFunc('/json', (w, r) => {
    w.header['Content-Type'] = 'application/json'
    w.write(JSON.stringify({ message: 'ok' }))
})

console.log('Server starting on :8080')
http.listenAndServe(':8080', mux)
```
