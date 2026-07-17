# net/url

## 模块内容

```typescript
interface url {
    /**
     * 解析 URL 字符串
     * @param rawURL 原始 URL 字符串
     */
    parse(rawURL: string): [URL, Error];

    /**
     * 解析请求 URI 字符串（不支持相对 URL）
     * @param rawURL 原始 URI 字符串
     */
    parseRequestURI(rawURL: string): [URL, Error];

    /**
     * 解析 URL 查询字符串为 Values 对象
     * @param query 查询字符串，如 "a=1&b=2"
     */
    parseQuery(query: string): [Values, Error];

    /**
     * 对字符串进行 URL 编码（百分号转义）
     * @param s 待编码的字符串
     */
    queryEscape(s: string): string;

    /**
     * 对 URL 编码的字符串进行解码
     * @param s 已编码的字符串
     */
    queryUnescape(s: string): [string, Error];

    /**
     * 创建仅包含用户名的 Userinfo
     * @param username 用户名
     */
    user(username: string): Userinfo;

    /**
     * 创建包含用户名和密码的 Userinfo
     * @param username 用户名
     * @param password 密码
     */
    userPassword(username: string, password: string): Userinfo;

    /** Values 构造函数，用于创建查询参数对象 */
    Values: {
        (): Values;
    };
}
```

### URL

```typescript
/**
 * 表示一个解析后的 URL
 */
interface URL {
    /** 协议，如 "https" */
    readonly scheme: string;
    /** 不透明部分（用于 mailto: 等非分层 URL） */
    readonly opaque: string;
    /** 主机和端口，如 "example.com:8080" */
    readonly host: string;
    /** 路径部分 */
    readonly path: string;
    /** 编码后的路径 */
    readonly rawPath: string;
    /** 编码后的查询字符串 */
    readonly rawQuery: string;
    /** URL 片段（# 后的部分） */
    readonly fragment: string;
    /** 编码后的片段 */
    readonly rawFragment: string;
    /** 是否在查询字符串为空时仍追加 "?" */
    forceQuery: boolean;
    /** 是否省略 URL 字符串中的 host 部分 */
    omitHost: boolean;

    /** 返回重新组装的完整 URL 字符串 */
    string(): string;

    /** 返回用于 HTTP 请求的 "path?query" 或 "opaque?query" 格式 */
    requestURI(): string;

    /** 将 ref 作为相对引用解析，返回新的 URL */
    resolveReference(ref: URL): URL;

    /** 解析 rawQuery 为 Values 对象 */
    query(): Values;

    /** 返回隐藏密码的 URL 字符串（密码替换为 "xxxxx"） */
    redacted(): string;

    /** 是否为绝对 URL（有 scheme） */
    isAbs(): boolean;

    /** 返回主机名部分（不含端口） */
    hostname(): string;

    /** 返回端口部分，没有端口返回空字符串 */
    port(): string;

    /** 将 elem 拼接到路径上并返回新 URL */
    joinPath(elem: string): URL;
}
```

### Values

```typescript
/**
 * URL 查询参数容器，一个键可对应多个值
 */
interface Values {
    /** 为指定键追加一个值（同键已有值则追加） */
    add(key: string, value: string): void;

    /** 删除指定键的所有值 */
    del(key: string): void;

    /** 获取指定键的第一个值，不存在返回空字符串 */
    get(key: string): string;

    /** 获取指定键的所有值 */
    getAll(key: string): string[];

    /** 判断是否存在指定键 */
    has(key: string): boolean;

    /** 设置指定键的值（覆盖已有值） */
    set(key: string, value: string): void;

    /** 将所有参数编码为 "key=value&..." 格式的查询字符串 */
    encode(): string;
}
```

### Userinfo

```typescript
/**
 * URL 中的用户认证信息
 */
interface Userinfo {
    /** 获取用户名 */
    username(): string;

    /** 获取密码，ok 表示是否设置了密码 */
    password(): [string, boolean];

    /** 返回 "username:password" 格式字符串 */
    string(): string;
}
```

## 使用样例

```javascript
import console from 'console'
import url from 'net/url'

// parse URL
const u = url.parse('https://example.com/path?q=hello&lang=zh#top')
console.log(u[0].scheme)       // "https"
console.log(u[0].host)         // "example.com"
console.log(u[0].path)         // "/path"
console.log(u[0].fragment)     // "top"
console.log(u[0].query().get('q')) // "hello"

// build query string
const values = new url.Values()
values.set('page', '1')
values.set('size', '20')
console.log(values.encode())   // "page=1&size=20"

// parse query string
const [q, err] = url.parseQuery('a=1&b=2&a=3')
console.log(q.getAll('a'))     // ["1", "3"]
console.log(q.get('b'))        // "2"
```
