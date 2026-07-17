# io

## 模块内容

```typescript
interface io {
    /**
     * 从 src 读取数据写入 dst，直到 EOF 或出错，返回写入字节数
     * @param dst 目标写入器
     * @param src 源读取器
     */
    copy(dst: Writer, src: Reader): [bigint, Error];

    /** 丢弃所有写入数据的 Writer（类似 /dev/null） */
    readonly discard: Writer;

    /** 表示数据流结束的哨兵错误，read 方法在流末尾返回此错误 */
    readonly EOF: Error;

    /**
     * 从 r 包装一个限Reader, 最多只读取 n 字节
     * @param r 源读取器
     * @param n 最大读取字节数
     */
    limitReader(r: Reader, n: bigint): Reader;

    /**
     * 从 r 读取全部数据直到 EOF，返回完整内容
     * @param r 源读取器
     */
    readAll(r: Reader): [Uint8Array, Error];
}
```

### Reader

```typescript
/**
 * 读取器接口，表示可读取字节流的数据源
 */
interface Reader {
    /**
     * 读取数据到缓冲区 p，返回实际读取的字节数
     * 流结束时返回 io.EOF 错误
     */
    read(p: Uint8Array): [number, Error];
}
```

### Writer

```typescript
/**
 * 写入器接口，表示可写入字节流的数据目标
 */
interface Writer {
    /** 将数据 p 写入，返回实际写入的字节数 */
    write(p: Uint8Array): [number, Error];
}
```

## 使用样例

```javascript
import console from 'console'
import io from 'io'
import http from 'net/http'

const res = http.fetch('https://httpbin.io/get')
const body = io.readAll(res.body)
console.log(body[0])
```
