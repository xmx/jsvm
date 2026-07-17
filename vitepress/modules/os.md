# os

## 模块内容

```typescript
interface os {
    /** 命令行参数列表，args[0] 为程序名称 */
    readonly args: string[];

    /** 返回当前可执行文件的路径 */
    executable(): [string, Error];

    /** 返回所有环境变量，格式为 "key=value" */
    environ(): string[];

    /** 返回主机名 */
    hostname(): [string, Error];

    /** 返回用户配置目录（如 Linux 下 ~/.config） */
    userConfigDir(): [string, Error];

    /** 返回用户缓存目录（如 Linux 下 ~/.cache） */
    userCacheDir(): [string, Error];

    /** 返回用户主目录 */
    userHomeDir(): [string, Error];

    /**
     * 获取指定环境变量的值，不存在返回空字符串
     * @param key 环境变量名
     */
    getenv(key: string): string;

    /** 返回当前进程的有效用户 ID（Unix） */
    geteuid(): number;

    /** 返回当前进程的组 ID */
    getgid(): number;

    /** 返回父进程的 PID */
    getppid(): number;

    /** 返回当前进程的 PID */
    getpid(): number;

    /** 返回当前工作目录 */
    getwd(): [string, Error];
}
```

## 使用样例

```javascript
import console from 'console'
import os from 'os'

console.log('PID:', os.getpid())
console.log('Hostname:', os.hostname())
console.log('Home:', os.userHomeDir())
console.log('CWD:', os.getwd())
console.log('HOME env:', os.getenv('HOME'))
```
