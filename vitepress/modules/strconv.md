# strconv

## 模块内容

```typescript
interface strconv {
    /**
     * 将字符串解析为有符号整数
     * @param s 待解析字符串
     * @param base 进制（2-36），0 表示自动推断（0x 前缀为 16 进制等）
     * @param bitSize 位宽（0、8、16、32、64）
     */
    parseInt(s: string, base: number, bitSize: number): [bigint, Error];

    /**
     * 将字符串解析为无符号整数
     * @param s 待解析字符串
     * @param base 进制（2-36）
     * @param bitSize 位宽（0、8、16、32、64）
     */
    parseUint(s: string, base: number, bitSize: number): [bigint, Error];

    /**
     * 将字符串解析为浮点数
     * @param s 待解析字符串
     * @param bitSize 32 表示 float32，64 表示 float64
     */
    parseFloat(s: string, bitSize: number): [number, Error];

    /**
     * 将字符串解析为布尔值，接受 "1"、"t"、"T"、"TRUE"、"true" 等为 true
     * @param str 待解析字符串
     */
    parseBool(str: string): [boolean, Error];

    /**
     * 将有符号整数格式化为指定进制的字符串
     * @param i 整数
     * @param base 进制（2-36）
     */
    formatInt(i: bigint, base: number): string;

    /**
     * 将无符号整数格式化为指定进制的字符串
     * @param i 无符号整数
     * @param base 进制（2-36）
     */
    formatUint(i: bigint, base: number): string;

    /**
     * 将浮点数格式化为字符串
     * @param f 浮点数
     * @param fmt 格式字符：'b'(-ddd.dddp±dd)、'e'(-d.dddde±dd)、'f'(-ddd.dddd)、'g'(自动选择)
     * @param prec 精度（小数点后位数，-1 表示最小位数）
     * @param bitSize 32 或 64
     */
    formatFloat(f: number, fmt: string, prec: number, bitSize: number): string;

    /** 将布尔值格式化为 "true" 或 "false" */
    formatBool(b: boolean): string;

    /** 将整数转为十进制字符串（等价于 formatInt(int64(i), 10)） */
    itoa(i: number): string;

    /** 将字符串解析为十进制整数（等价于 parseInt(s, 10, 0) 返回 int） */
    atoi(s: string): [number, Error];

    /** 将字符串用双引号包裹为 Go 字符串字面量形式 */
    quote(s: string): string;

    /** 将字符串用双引号包裹并转义非 ASCII 字符为 \uXXXX 形式 */
    quoteToASCII(s: string): string;

    /**
     * 解析 Go 字符串字面量（去除外层引号和转义）
     * @param s 带引号的字符串，如 `"hello"` 或 `'hello'`
     */
    unquote(s: string): [string, Error];

    /** 判断 rune 是否可以用反引号字符串字面量表示（不含 '`' 和 '\\'） */
    canBackquote(r: number): boolean;

    /** 判断 rune 是否为可打印字符 */
    isPrint(r: number): boolean;

    /** 判断 rune 是否为可打印或空格字符 */
    isGraphic(r: number): boolean;
}
```

## 使用样例

```javascript
import console from 'console'
import strconv from 'strconv'

// int to string
console.log(strconv.itoa(42))           // "42"
console.log(strconv.atoi("42"))          // [42, null]
console.log(strconv.formatBool(true))    // "true"
console.log(strconv.parseBool("true"))   // [true, null]

// quoting
console.log(strconv.quote("hello"))      // `"hello"`
console.log(strconv.unquote('"hello"'))  // ["hello", null]
```
