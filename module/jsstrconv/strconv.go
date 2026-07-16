package jsstrconv

import (
	"strconv"

	"github.com/xmx/jsvm"
)

type strconvModule struct{}

func New() jsvm.ModuleExporter {
	return &strconvModule{}
}

func (m *strconvModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	vals := map[string]any{
		"parseInt":     strconv.ParseInt,
		"parseUint":    strconv.ParseUint,
		"parseFloat":   strconv.ParseFloat,
		"parseBool":    strconv.ParseBool,
		"formatInt":    strconv.FormatInt,
		"formatUint":   strconv.FormatUint,
		"formatFloat":  strconv.FormatFloat,
		"formatBool":   strconv.FormatBool,
		"itoa":         strconv.Itoa,
		"atoi":         strconv.Atoi,
		"quote":        strconv.Quote,
		"quoteToASCII": strconv.QuoteToASCII,
		"unquote":      strconv.Unquote,
		"canBackquote": strconv.CanBackquote,
		"isPrint":      strconv.IsPrint,
		"isGraphic":    strconv.IsGraphic,
	}

	return jsvm.ModuleExports{
		Name:    "strconv",
		Default: vals,
	}
}
