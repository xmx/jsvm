package jsstrconv

import (
	"strconv"

	"github.com/xmx/jsvm"
)

// strconvModule 是 strconv 模块的实现，将 Go 字符串/数字转换函数导出到 JS。
type strconvModule struct{}

// New 创建 strconv 模块实例。
func New() jsvm.ModuleExporter {
	return &strconvModule{}
}

// ModuleExports 注册 strconv 模块，导出：
//
// 解析类：parseInt/parseUint/parseFloat/parseBool — 字符串转数字/布尔
// 格式化：formatInt/formatUint/formatFloat/formatBool — 数字/布尔转字符串
// 十进制：itoa/atoi — 十进制转换快捷函数
// 字符串字面量：quote/quoteToASCII/unquote — 引号包裹与解包
// 字符分类：canBackquote/isPrint/isGraphic — rune 属性判断
func (m *strconvModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	vals := map[string]any{
		"parseInt":     strconv.ParseInt,     // 字符串转有符号整数（指定进制和位宽）
		"parseUint":    strconv.ParseUint,    // 字符串转无符号整数
		"parseFloat":   strconv.ParseFloat,   // 字符串转浮点数
		"parseBool":    strconv.ParseBool,    // 字符串转布尔（接受 "1"/"true"/"TRUE" 等）
		"formatInt":    strconv.FormatInt,    // 有符号整数转指定进制字符串
		"formatUint":   strconv.FormatUint,   // 无符号整数转指定进制字符串
		"formatFloat":  strconv.FormatFloat,  // 浮点数转字符串（指定格式和精度）
		"formatBool":   strconv.FormatBool,   // 布尔转 "true"/"false"
		"itoa":         strconv.Itoa,         // 十进制整数转字符串（formatInt(i, 10) 的快捷形式）
		"atoi":         strconv.Atoi,         // 字符串转十进制整数
		"quote":        strconv.Quote,        // 将字符串用双引号包裹，转义非 ASCII 字符
		"quoteToASCII": strconv.QuoteToASCII, // 将字符串用双引号包裹并转义非 ASCII 为 \uXXXX
		"unquote":      strconv.Unquote,      // 解析带引号的字符串字面量（去除引号和转义）
		"canBackquote": strconv.CanBackquote, // 判断 rune 是否可用于反引号字符串
		"isPrint":      strconv.IsPrint,      // 判断 rune 是否为可打印字符
		"isGraphic":    strconv.IsGraphic,    // 判断 rune 是否为可见字符（含空格）
	}

	return jsvm.ModuleExports{
		Name:    "strconv",
		Default: vals,
	}
}
