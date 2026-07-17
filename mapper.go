package jsvm

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/dop251/goja/parser"
)

// tagMapper 是 goja 字段名映射器，将 Go struct 字段/方法名转为 JS 风格。
// 使用 struct tag 中指定名称，无 tag 时使用驼峰规则降序。
// 例如：传入 "json" 表示使用 json tag。
type tagMapper string

// FieldName 返回 struct 字段在 JS 中使用的名称。
// 优先使用 tag 值（去掉逗号后的部分），tag 无效时回退到降序驼峰转换。
func (tm tagMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	tag := f.Tag.Get(string(tm))
	// 去掉 tag 中逗号及以后的选项部分（如 `json:"name,omitempty"` 取 "name"）
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	// tag 必须是合法的 JS 标识符才使用，否则使用降低首字母后的字段名
	if parser.IsIdentifier(tag) {
		return tag
	}

	return tm.lowerCase(f.Name)
}

// MethodName 返回 struct 方法在 JS 中使用的名称。
// 将 Go 导出方法名（大写首字母）转为 JS 风格（小写首字母）。
func (tm tagMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return tm.lowerCase(m.Name)
}

// lowerCase 将 Go 导出名称转为 JS 风格（首字母小写，连续大写按单词拆分）。
//
// 转换规则：
//
//	HTTP       -> http        （全大写字段直接全部小写）
//	MyHTTP     -> myHTTP      （前缀单词小写，后续缩写保持大写）
//	CopyN      -> copyN       （末尾大写字母保留）
//	N          -> n
//	JSONBody   -> jsonBody    （开头连续大写全部小写，直到遇到小写字母）
//
// 逻辑：从前往后遍历，将开头连续大写字母（以及紧邻小写字母之前的大写）都转为小写，
// 一旦遇到小写字母或大写字母后面紧跟小写字母（词边界），停止转换。
func (tagMapper) lowerCase(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	for i, r := range runes {
		// 遇到小写字母直接停止转换
		if unicode.IsLower(r) {
			break
		}
		// 只要不是首字母，且下一个字符是小写字母，说明当前大写字母是新单词的开头，停止转换。
		// 例 1: XMLName -> 遍历到 N 时，下一个是 a，N 保持大写
		// 例 2: HTML5String -> 遍历到 S 时，下一个是 t，S 保持大写
		if i > 0 && i < len(runes)-1 && unicode.IsLower(runes[i+1]) {
			break
		}
		runes[i] = unicode.ToLower(r)
	}

	return string(runes)
}
