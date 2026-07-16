package jsstrings

import (
	"strings"

	"github.com/xmx/jsvm"
)

type stringsModule struct{}

func New() jsvm.ModuleExporter {
	return &stringsModule{}
}

func (m *stringsModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	vals := map[string]any{
		"contains":     strings.Contains,
		"containsAny":  strings.ContainsAny,
		"hasPrefix":    strings.HasPrefix,
		"hasSuffix":    strings.HasSuffix,
		"index":        strings.Index,
		"indexAny":     strings.IndexAny,
		"lastIndex":    strings.LastIndex,
		"lastIndexAny": strings.LastIndexAny,
		"count":        strings.Count,
		"repeat":       strings.Repeat,
		"replace":      strings.Replace,
		"replaceAll":   strings.ReplaceAll,
		"split":        strings.Split,
		"splitN":       strings.SplitN,
		"splitAfter":   strings.SplitAfter,
		"splitAfterN":  strings.SplitAfterN,
		"fields":       strings.Fields,
		"join":         strings.Join,
		"toLower":      strings.ToLower,
		"toUpper":      strings.ToUpper,
		"toTitle":      strings.ToTitle,
		"toValidUTF8":  strings.ToValidUTF8,
		"trim":         strings.Trim,
		"trimLeft":     strings.TrimLeft,
		"trimRight":    strings.TrimRight,
		"trimSpace":    strings.TrimSpace,
		"trimPrefix":   strings.TrimPrefix,
		"trimSuffix":   strings.TrimSuffix,
		"compare":      strings.Compare,
		"equalFold":    strings.EqualFold,
		"newReader":    strings.NewReader,
		"newReplacer":  strings.NewReplacer,
	}

	return jsvm.ModuleExports{
		Name:    "strings",
		Default: vals,
	}
}
