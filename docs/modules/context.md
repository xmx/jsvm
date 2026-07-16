# context

## 模块内容

```typescript
interface context {
    /** 返回一个永不过期的空 Context，通常作为根 Context 使用 */
    background(): Context;

    /**
     * 基于父 Context 创建可取消的子 Context
     * @returns [子 Context, 取消函数]
     */
    withCancel(parent: Context): [Context, CancelFunc];

    /**
     * 基于父 Context 创建带超时的子 Context
     * @param d 超时时长
     * @returns [子 Context, 取消函数]
     */
    withTimeout(parent: Context, d: Duration): [Context, CancelFunc];

    /**
     * 基于父 Context 创建带截止时间的子 Context
     * @param d 截止时间
     * @returns [子 Context, 取消函数]
     */
    withDeadline(parent: Context, d: Date): [Context, CancelFunc];

    /**
     * 基于父 Context 创建携带键值对的子 Context
     * @param key 键
     * @param val 值
     */
    withValue(parent: Context, key: any, val: any): Context;
}
```

### CancelFunc

```typescript
/** 取消函数，调用后取消关联的 Context */
type CancelFunc = () => void;
```

### Context

```typescript
/**
 * 上下文对象，用于传递截止时间、取消信号和请求级别的值
 */
interface Context {
    /** 返回此 Context 的截止时间，如果没有截止时间则 ok 为 false */
    deadline(): [Date, boolean];

    /** 返回 Context 被取消时的错误，未取消返回 null */
    err(): Error;

    /** 获取与指定键关联的值，不存在返回 undefined */
    value(key: any): any;
}
```

## 使用样例

```javascript
import console from 'console'
import context from 'context'
import time from 'time'

// 创建一个可取消的 Context
const [ctx, cancel] = context.withCancel(context.background())
// ... 使用 ctx ...
cancel()

// 创建一个带超时的 Context
const [ctx2, cancel2] = context.withTimeout(context.background(), 5 * time.second)
delay(5 * time.second, 'done').then(v => {
    console.log(v)
    cancel2()
})

// 传递值
let ctx3 = context.withValue(context.background(), 'key', 'val')
console.log(ctx3.value('key'))  // "val"
```
