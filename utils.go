package jsvm

import (
	"encoding/base64"
	"encoding/json"

	"github.com/dop251/goja"
)

// StringValue 将 goja 值转换为可读的字符串表示。
// 处理规则：
//   - 构造函数返回 "<Constructor Function>"
//   - 普通函数返回 "<Function>"
//   - Error 对象返回 "ErrorType: message"
//   - []byte 类型直接转为字符串
//   - 其他对象使用 JSON 序列化
//   - Symbol 返回 "Symbol(value)"
//   - 其他值直接调用 toString()
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
			// 特殊处理 Error 类型，展示错误类型和消息
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
			// 字节数组直接转为字符串
			return string(gv)
		default:
			// 其他类型用 JSON 序列化
			bs, err := json.Marshal(ev)
			if err != nil {
				return "Marshal Error: " + err.Error()
			}

			return string(bs)
		}
	case *goja.Symbol:
		// Symbol 按 JS 规范格式输出
		s := tv.String()
		return "Symbol(" + s + ")"
	}

	return v.String()
}

// StringArrayBuffer 将 ArrayBuffer 的字节内容编码为 Base64 字符串。
func StringArrayBuffer(ab goja.ArrayBuffer) string {
	bs := ab.Bytes()
	return base64.StdEncoding.EncodeToString(bs)
}

// IsNullish 判断值是否为空值（nil、undefined 或 null）。
// 这三个值在 JavaScript 中均视为"不存在"，常用于参数校验。
func IsNullish(v goja.Value) bool {
	return v == nil || goja.IsNull(v) || goja.IsUndefined(v)
}
