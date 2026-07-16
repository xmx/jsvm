# net/http

## 模块内容

```typescript
interface http {
    fetch(url: string, info?: RequestInfo): Response
}
```

### RequestInfo

```typescript
interface RequestInfo {
    method: string;
    header: { [key: string]: string };
    body: string;
}
```

### Response

```typescript
interface Response {
    readonly status: string;
    readonly statusCode: number;

    text(): string;

    json(): any;
}
```

## 使用样例

```javascript
import console from 'console'
import http from 'net/http'
import url from 'net/url'

const params = new url.Values()
params.set('host', '127.0.0.1')
params.set('port', 443)

const res = http.fetch('https://httpbin.io/get?' + params.encode(), {
    method: 'GET',
    header: {
        Accept: 'application/json',
        Host: 'localhost',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36 Edg/150.0.0.0',
    },
    body: params,
});

const data = res.text()
console.log(data);
```
