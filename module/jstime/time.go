package jstime

import (
	"time"

	"github.com/xmx/jsvm"
)

func New() jsvm.ModuleExporter {
	return &timeModule{}
}

type timeModule struct{}

func (m *timeModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		"nanosecond":    time.Nanosecond,
		"microsecond":   time.Microsecond,
		"millisecond":   time.Millisecond,
		"second":        time.Second,
		"minute":        time.Minute,
		"hour":          time.Hour,
		"january":       time.January,
		"february":      time.February,
		"march":         time.March,
		"april":         time.April,
		"may":           time.May,
		"june":          time.June,
		"july":          time.July,
		"august":        time.August,
		"september":     time.September,
		"october":       time.October,
		"november":      time.November,
		"december":      time.December,
		"sunday":        time.Sunday,
		"monday":        time.Monday,
		"tuesday":       time.Tuesday,
		"wednesday":     time.Wednesday,
		"thursday":      time.Thursday,
		"friday":        time.Friday,
		"saturday":      time.Saturday,
		"local":         time.Local,
		"parseDuration": time.ParseDuration,
		"loadLocation":  time.LoadLocation,
		"date":          time.Date,
	}

	return jsvm.ModuleExports{
		Name:    "time",
		Default: defaults,
	}
}
