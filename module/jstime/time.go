package jstime

import (
	"time"

	"github.com/xmx/jsvm"
)

// New 创建 time 模块实例。
func New() jsvm.ModuleExporter {
	return &timeModule{}
}

// timeModule 是 time 模块的实现，将 Go time 包的常量、月份、星期及函数导出到 JS。
type timeModule struct{}

// ModuleExports 注册 time 模块，导出：
//   - Duration 常量：nanosecond/microsecond/millisecond/second/minute/hour
//   - Month 常量：january 到 december
//   - Weekday 常量：sunday 到 saturday
//   - 时区：local（系统本地时区）
//   - 函数：parseDuration（解析时长字符串）、loadLocation（加载时区）、date（构造时间）
func (m *timeModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		// Duration 常量
		"nanosecond":  time.Nanosecond,
		"microsecond": time.Microsecond,
		"millisecond": time.Millisecond,
		"second":      time.Second,
		"minute":      time.Minute,
		"hour":        time.Hour,
		// Month 常量
		"january":   time.January,
		"february":  time.February,
		"march":     time.March,
		"april":     time.April,
		"may":       time.May,
		"june":      time.June,
		"july":      time.July,
		"august":    time.August,
		"september": time.September,
		"october":   time.October,
		"november":  time.November,
		"december":  time.December,
		// Weekday 常量
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
		// 本地时区
		"local": time.Local,
		// 函数
		"parseDuration": time.ParseDuration, // 解析 "1h30m" 格式的时长字符串
		"loadLocation":  time.LoadLocation,  // 按 IANA 时区名加载 Location
		"date":          time.Date,          // 构造指定时间点的 Time 对象
	}

	return jsvm.ModuleExports{
		Name:    "time",
		Default: defaults,
	}
}
