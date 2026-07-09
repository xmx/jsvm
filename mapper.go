package jsvm

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/dop251/goja/parser"
)

type tagMapper string

func (tm tagMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	tag := f.Tag.Get(string(tm))
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	if parser.IsIdentifier(tag) {
		return tag
	}

	return tm.lowerCase(f.Name)
}

func (tm tagMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return tm.lowerCase(m.Name)
}

// lowerCase 将 Go 可导出变量转为 JS 风格的变量。
//
//	HTTP -> http
//	MyHTTP -> myHTTP
//	CopyN -> copyN
//	N -> n
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
