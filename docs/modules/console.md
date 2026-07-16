# console

## 模块内容

```typescript
interface console {

    /**
     * 等价于 @see{info} 方法
     *
     * @param data
     */
    log(...data: any[]): void;

    debug(...data: any[]): void;

    info(...data: any[]): void;

    warn(...data: any[]): void;

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
