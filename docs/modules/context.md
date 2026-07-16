# context

## 模块内容

```typescript
interface context {
    background(): Context;

    withCancel(parent: Context): [Context, CancelFunc];
}
```

### CancelFunc

```typescript
type CancelFunc = () => void;
```

### Context

```typescript
interface Context {
    deadline(): [Date, boolean];

    err(): Error;

    value(key: any): any;
}
```

## 使用样例

```javascript
import console from 'console'

console.log('hello')
console.info('hello')
console.error('hello')
```
