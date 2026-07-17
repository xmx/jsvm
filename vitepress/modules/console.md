# console

## 模块内容

```typescript
interface console {
    /**
     * 输出信息日志，等价于 info 方法
     * @param data 要输出的数据
     */
    log(...data: any[]): void;

    /**
     * 输出调试级别的日志，需要日志级别设为 Debug 才会实际输出
     * @param data 要输出的数据
     */
    debug(...data: any[]): void;

    /**
     * 输出信息级别的日志
     * @param data 要输出的数据
     */
    info(...data: any[]): void;

    /**
     * 输出警告级别的日志
     * @param data 要输出的数据
     */
    warn(...data: any[]): void;

    /**
     * 输出错误级别的日志
     * @param data 要输出的数据
     */
    error(...data: any[]): void;
}
```

## 使用样例

```javascript
import console from 'console'

console.log('hello')
console.info('hello')
console.error('hello')
```
