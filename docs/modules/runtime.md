# runtime

## 模块内容

```typescript
interface runtime {
    /** 读取内存使用统计信息（会触发 GC） */
    readMemStats(): MemStats;

    /** 当前操作系统，如 "linux"、"windows"、"darwin" */
    readonly GOOS: string;

    /** 当前 CPU 架构，如 "amd64"、"arm64" */
    readonly GOARCH: string;

    /** 编译器名称，如 "gc" */
    readonly Compiler: string;

    /** 返回 cgo 调用次数 */
    numCgoCall(): bigint;

    /** 返回可用的 CPU 核心数 */
    numCPU(): number;

    /** 返回当前存活的 goroutine 数量 */
    numGoroutine(): number;

    /** 返回 Go 版本号，如 "go1.22.0" */
    version(): string;
}
```

### MemStats

```typescript
/**
 * 内存统计信息，由 readMemStats 返回
 */
interface MemStats {
    /** 已分配且存活的堆内存字节数 */
    readonly alloc: bigint;
    /** 生命周期内累计分配的总堆内存字节数 */
    readonly totalAlloc: bigint;
    /** 从操作系统获取的总内存字节数 */
    readonly sys: bigint;
    /** 累计堆分配次数 */
    readonly mallocs: bigint;
    /** 累计释放次数 */
    readonly frees: bigint;
    /** 已分配且存活的堆内存字节数（同 alloc） */
    readonly heapAlloc: bigint;
    /** 从操作系统获取的堆内存字节数 */
    readonly heapSys: bigint;
    /** 空闲的堆内存字节数 */
    readonly heapIdle: bigint;
    /** 正在使用的堆内存字节数 */
    readonly heapInuse: bigint;
    /** 归还给操作系统的堆内存字节数 */
    readonly heapReleased: bigint;
    /** 存活的堆对象数量 */
    readonly heapObjects: bigint;
    /** 正在使用的栈内存字节数 */
    readonly stackInuse: bigint;
    /** 从操作系统获取的栈内存字节数 */
    readonly stackSys: bigint;
    /** 累计 GC 次数 */
    readonly numGC: bigint;
    /** 累计强制 GC 次数 */
    readonly numForcedGC: bigint;
    /** 上次 GC 完成的 Unix 纳秒时间戳 */
    readonly lastGC: bigint;
    /** 程序启动以来 GC 暂停累计纳秒数 */
    readonly pauseTotalNs: bigint;
}
```

## 使用样例

```javascript
import console from 'console'
import runtime from 'runtime'

console.log('Go version:', runtime.version())
console.log('OS/Arch:', runtime.GOOS, runtime.GOARCH)
console.log('CPU:', runtime.numCPU())
console.log('Goroutines:', runtime.numGoroutine())

const stats = runtime.readMemStats()
console.log('Heap alloc:', stats.heapAlloc)
```
