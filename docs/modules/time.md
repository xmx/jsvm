# time

## 模块内容

```typescript
interface time {
    nanosecond: Duration;
    microsecond: Duration;
    millisecond: Duration;
    second: Duration;
    minute: Duration;
    hour: Duration;

    // "january":       time.January,
    // "february":      time.February,
    // "march":         time.March,
    // "april":         time.April,
    // "may":           time.May,
    // "june":          time.June,
    // "july":          time.July,
    // "august":        time.August,
    // "september":     time.September,
    // "october":       time.October,
    // "november":      time.November,
    // "december":      time.December,
    sunday: Weekday,
    monday: Weekday,
    tuesday: Weekday,
    wednesday: Weekday,
    thursday: Weekday,
    friday: Weekday,
    saturday: Weekday,

    // "local":         time.Local,
    parseDuration(str: string): Duration;
}
```

### Duration

```typescript
interface Duration {
    string(): string;

    nanoseconds(): bigint;
    
    microseconds(): bigint;
    
    milliseconds(): bigint;

    seconds(): number;

    minutes(): number;

    hours(): number;

    truncate(d: Duration): Duration;

    round(d: Duration): Duration;
    
    abs(): Duration;
}
```

### Weekday

```typescript
interface Weekday {
    string(): string;
}
```

### Month

```typescript
interface Month {
    string(): string;
}
```

## 使用样例

```javascript

```
