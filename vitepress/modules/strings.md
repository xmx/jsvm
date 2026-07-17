# strings

## 模块内容

```typescript
interface strings {
    /** 判断 s 是否包含子串 substr */
    contains(s, substr: string): boolean;

    /** 判断 s 是否包含 chars 中的任意字符 */
    containsAny(s, chars: string): boolean;

    /** 判断 s 是否以 prefix 开头 */
    hasPrefix(s, prefix: string): boolean;

    /** 判断 s 是否以 suffix 结尾 */
    hasSuffix(s, suffix: string): boolean;

    /** 返回 substr 在 s 中首次出现的索引，未找到返回 -1 */
    index(s, substr: string): number;

    /** 返回 chars 中任意字符在 s 中首次出现的索引，未找到返回 -1 */
    indexAny(s, chars: string): number;

    /** 返回 substr 在 s 中最后一次出现的索引，未找到返回 -1 */
    lastIndex(s, substr: string): number;

    /** 返回 chars 中任意字符在 s 中最后一次出现的索引，未找到返回 -1 */
    lastIndexAny(s, chars: string): number;

    /** 返回 substr 在 s 中出现的次数 */
    count(s, substr: string): number;

    /**
     * 将 s 重复 count 次并拼接返回
     * @param count 重复次数，0 返回空字符串
     */
    repeat(s: string, count: number): string;

    /**
     * 替换 s 中的 old 为 new，最多替换 n 次；n < 0 表示全部替换
     */
    replace(s, old, new: string, n: number): string;

    /** 替换 s 中所有 old 为 new */
    replaceAll(s, old, new: string): string;

    /** 按 sep 分割 s，返回所有子串数组 */
    split(s, sep: string): string[];

    /**
     * 按 sep 分割 s，最多返回 n 个子串
     * @param n 最大子串数
     */
    splitN(s, sep: string, n: number): string[];

    /** 按 sep 分割 s，每个子串保留 sep 后缀 */
    splitAfter(s, sep: string): string[];

    /** 按 sep 分割 s，每个子串保留 sep 后缀，最多返回 n 个子串 */
    splitAfterN(s, sep: string, n: number): string[];

    /** 按空白字符分割 s，返回非空白子串数组 */
    fields(s: string): string[];

    /** 用 sep 将字符串数组拼接为一个字符串 */
    join(elems: string[], sep: string): string;

    /** 将所有字母转为小写 */
    toLower(s: string): string;

    /** 将所有字母转为大写 */
    toUpper(s: string): string;

    /** 将字符串转为标题形式（每个 Unicode 字母词首大写） */
    toTitle(s: string): string;

    /**
     * 将无效 UTF-8 字节替换为指定字符串
     * @param replacement 替换用的字符串
     */
    toValidUTF8(s, replacement: string): string;

    /** 去除 s 两端属于 cutset 的字符 */
    trim(s, cutset: string): string;

    /** 去除 s 左侧属于 cutset 的字符 */
    trimLeft(s, cutset: string): string;

    /** 去除 s 右侧属于 cutset 的字符 */
    trimRight(s, cutset: string): string;

    /** 去除 s 两端的空白字符（空格、制表符、换行等） */
    trimSpace(s: string): string;

    /** 去除 s 的 prefix 前缀，不匹配则返回原字符串 */
    trimPrefix(s, prefix: string): string;

    /** 去除 s 的 suffix 后缀，不匹配则返回原字符串 */
    trimSuffix(s, suffix: string): string;

    /**
     * 按字典序比较两个字符串
     * @returns 0 相等，-1 a < b，+1 a > b
     */
    compare(a, b: string): number;

    /** 忽略大小写比较两个字符串是否相等 */
    equalFold(a, b: string): boolean;

    /** 将字符串包装为一个 Reader */
    newReader(s: string): Reader;

    /**
     * 创建多对替换器，参数为 old/new 交替的字符串对
     * @example newReplacer("old1", "new1", "old2", "new2")
     */
    newReplacer(...oldnew: string[]): Replacer;
}
```

### Reader

```typescript
/**
 * 字符串 Reader，实现了 io.Reader 接口
 */
interface Reader {
    /** 读取数据到缓冲区 p，返回实际读取的字节数 */
    read(p: Uint8Array): [number, Error];
}
```

### Replacer

```typescript
/**
 * 多对替换器，由 newReplacer 创建
 */
interface Replacer {
    /** 对 s 执行所有替换并返回结果 */
    replace(s: string): string;

    /** 将替换后的内容写入指定 Writer，返回写入字节数 */
    writeString(w: Writer, s: string): [number, Error];
}
```

## 使用样例

```javascript
import console from 'console'
import strings from 'strings'

// split & join
console.log(strings.split("a,b,c", ","))       // ["a","b","c"]
console.log(strings.join(["x", "y"], "-"))     // "x-y"

// search
console.log(strings.contains("hello", "lo"))   // true
console.log(strings.index("hello", "l"))       // 2

// case conversion
console.log(strings.toUpper("hello"))           // "HELLO"
console.log(strings.toLower("HELLO"))           // "hello"

// trim
console.log(strings.trimSpace("  hi  "))        // "hi"
console.log(strings.trimPrefix("GoGoGo", "Go")) // "GoGo"

// replace
console.log(strings.replaceAll("a-b-c", "-", "_")) // "a_b_c"
```
