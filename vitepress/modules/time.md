# time

## 模块内容

```typescript
interface time {
    /** 1 纳秒 */
    nanosecond: Duration;
    /** 1 微秒 = 1000 纳秒 */
    microsecond: Duration;
    /** 1 毫秒 = 1000 微秒 */
    millisecond: Duration;
    /** 1 秒 = 1000 毫秒 */
    second: Duration;
    /** 1 分钟 = 60 秒 */
    minute: Duration;
    /** 1 小时 = 60 分钟 */
    hour: Duration;

    /** 一月 */
    january: Month;
    /** 二月 */
    february: Month;
    /** 三月 */
    march: Month;
    /** 四月 */
    april: Month;
    /** 五月 */
    may: Month;
    /** 六月 */
    june: Month;
    /** 七月 */
    july: Month;
    /** 八月 */
    august: Month;
    /** 九月 */
    september: Month;
    /** 十月 */
    october: Month;
    /** 十一月 */
    november: Month;
    /** 十二月 */
    december: Month;

    /** 星期日 */
    sunday: Weekday;
    /** 星期一 */
    monday: Weekday;
    /** 星期二 */
    tuesday: Weekday;
    /** 星期三 */
    wednesday: Weekday;
    /** 星期四 */
    thursday: Weekday;
    /** 星期五 */
    friday: Weekday;
    /** 星期六 */
    saturday: Weekday;

    /** 本地时区 */
    readonly local: Location;

    /**
     * 解析时长字符串，如 "1h30m"、"500ms"
     * @param str 时长字符串
     */
    parseDuration(str: string): [Duration, Error];

    /**
     * 加载指定名称的时区，如 "Asia/Shanghai"、"UTC"
     * @param name 时区名称（IANA Time Zone Database 格式）
     */
    loadLocation(name: string): [Location, Error];

    /**
     * 创建指定时间点
     * @param year 年
     * @param month 月份（使用 time.january 等常量）
     * @param day 日
     * @param hour 时
     * @param min 分
     * @param sec 秒
     * @param nsec 纳秒
     * @param loc 时区
     */
    date(year: number, month: Month, day: number, hour: number, min: number, sec: number, nsec: number, loc: Location): Time;
}
```

### Duration

```typescript
/**
 * 时长类型，底层为纳秒数的 int64，支持算术运算
 */
interface Duration {
    /** 返回可读的时长字符串，如 "1h30m0s" */
    string(): string;

    /** 返回纳秒数 */
    nanoseconds(): bigint;

    /** 返回微秒数 */
    microseconds(): bigint;

    /** 返回毫秒数 */
    milliseconds(): bigint;

    /** 返回秒数（浮点） */
    seconds(): number;

    /** 返回分钟数（浮点） */
    minutes(): number;

    /** 返回小时数（浮点） */
    hours(): number;

    /**
     * 将时长截断到 d 的整数倍
     * @param d 截断粒度
     */
    truncate(d: Duration): Duration;

    /**
     * 将时长四舍五入到 d 的整数倍
     * @param d 四舍五入粒度
     */
    round(d: Duration): Duration;

    /** 返回时长的绝对值 */
    abs(): Duration;
}
```

### Location

```typescript
/** 时区 */
interface Location {
    /** 返回时区名称，如 "Asia/Shanghai" */
    string(): string;
}
```

### Weekday

```typescript
/** 星期几 */
interface Weekday {
    /** 返回星期名称，如 "Monday" */
    string(): string;
}
```

### Month

```typescript
/** 月份 */
interface Month {
    /** 返回月份名称，如 "January" */
    string(): string;
}
```

### Time

```typescript
/**
 * 时间点，精确到纳秒
 */
interface Time {
    /** 返回年份 */
    year(): number;
    /** 返回月份 */
    month(): Month;
    /** 返回日（1-31） */
    day(): number;
    /** 返回星期几 */
    weekday(): Weekday;
    /** 返回小时（0-23） */
    hour(): number;
    /** 返回分钟（0-59） */
    minute(): number;
    /** 返回秒（0-59） */
    second(): number;
    /** 返回纳秒部分（0-999999999） */
    nanosecond(): number;
    /** 是否为零值时间（January 1, year 1, 00:00:00 UTC） */
    isZero(): boolean;
    /** 返回 Unix 时间戳（秒） */
    unix(): bigint;
    /** 返回 Unix 时间戳（微秒） */
    unixMicro(): bigint;
    /** 返回 Unix 时间戳（毫秒） */
    unixMilli(): bigint;
    /** 返回 Unix 时间戳（纳秒） */
    unixNano(): bigint;
    /**
     * 按指定布局格式化时间字符串
     * @param layout 布局字符串（使用 Go 参考时间 Mon Jan 2 15:04:05 MST 2006）
     */
    format(layout: string): string;
    /** 加上指定时长，返回新的时间点 */
    add(d: Duration): Time;
    /** 计算与另一个时间的差值 */
    sub(u: Time): Duration;
    /** 是否与另一个时间相等 */
    equal(u: Time): boolean;
    /** 是否早于另一个时间 */
    before(u: Time): boolean;
    /** 是否晚于另一个时间 */
    after(u: Time): boolean;
    /** 返回当前时间所在的时区 */
    location(): Location;
    /**
     * 将时间转换到指定时区（同一时刻的不同表示）
     * @param loc 目标时区
     */
    in(loc: Location): Time;
}
```

## 使用样例

```javascript
import console from 'console'
import time from 'time'

// 解析时长
const [d, err] = time.parseDuration('1h30m')
console.log(d.string())   // "1h30m0s"
console.log(d.hours())    // 1.5

// 使用 Duration 常量
console.log((5 * time.second).string())   // "5s"

// 创建指定时间
const t = time.date(2025, time.january, 1, 12, 0, 0, 0, time.local)
console.log(t.year(), t.month(), t.day())  // 2025 January 1

// 加载时区
const [loc, _] = time.loadLocation('Asia/Shanghai')
const shanghai = t.in(loc)
console.log(shanghai.hour())
```
