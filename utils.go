package jsvm

import (
	"encoding/base64"
	"encoding/json"

	"github.com/dop251/goja"
)

func StringValue(v goja.Value) string {
	if _, ok := goja.AssertConstructor(v); ok {
		return "<Constructor Function>"
	}
	if _, ok := goja.AssertFunction(v); ok {
		return "<Function>"
	}
	switch tv := v.(type) {
	case *goja.Object:
		className := tv.ClassName()
		if className == "Error" {
			proto := tv.Prototype()
			errorType := proto.String()
			if errorType == "" {
				errorType = className
			}
			msg := tv.Get("message")
			var str string
			if msg != nil {
				str = msg.String()
			}

			return errorType + ": " + str
		}

		ev := tv.Export()
		switch gv := ev.(type) {
		case []byte:
			return string(gv)
		default:
			bs, err := json.Marshal(ev)
			if err != nil {
				return "Marshal Error: " + err.Error()
			}

			return string(bs)
		}
	case *goja.Symbol:
		s := tv.String()
		return "Symbol(" + s + ")"
	}

	return v.String()
}

func StringArrayBuffer(ab goja.ArrayBuffer) string {
	bs := ab.Bytes()
	return base64.StdEncoding.EncodeToString(bs)
}

// IsNullish checks if the given value is nullish, i.e. nil, undefined or null.
func IsNullish(v goja.Value) bool {
	return v == nil || goja.IsNull(v) || goja.IsUndefined(v)
}
