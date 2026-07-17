package jsstrings

import (
	"strings"

	"github.com/xmx/jsvm"
)

// stringsModule 是 strings 模块的实现，将 Go strings 包的常用函数导出到 JS。
type stringsModule struct{}

// New 创建 strings 模块实例。
func New() jsvm.ModuleExporter {
	return &stringsModule{}
}

// ModuleExports 注册 strings 模块，导出字符串查找、替换、分割、大小写转换、裁剪等函数。
func (m *stringsModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	vals := map[string]any{
		// 查找类
		"contains":     strings.Contains,     // 判断是否包含子串
		"containsAny":  strings.ContainsAny,  // 判断是否包含任意指定字符
		"hasPrefix":    strings.HasPrefix,    // 判断是否以指定前缀开头
		"hasSuffix":    strings.HasSuffix,    // 判断是否以指定后缀结尾
		"index":        strings.Index,        // 子串首次出现的索引
		"indexAny":     strings.IndexAny,     // 任意指定字符首次出现的索引
		"lastIndex":    strings.LastIndex,    // 子串最后一次出现的索引
		"lastIndexAny": strings.LastIndexAny, // 任意指定字符最后一次出现的索引
		"count":        strings.Count,        // 子串出现次数

		// 替换与重复
		"repeat":     strings.Repeat,     // 重复拼接字符串指定次数
		"replace":    strings.Replace,    // 替换（最多 n 次，n<0 为全部）
		"replaceAll": strings.ReplaceAll, // 替换所有

		// 分割
		"split":       strings.Split,       // 按分隔符分割为切片
		"splitN":      strings.SplitN,      // 按分隔符最多分割为 n 个子串
		"splitAfter":  strings.SplitAfter,  // 按分隔符分割并保留分隔符后缀
		"splitAfterN": strings.SplitAfterN, // 带数量限制的 splitAfter
		"fields":      strings.Fields,      // 按空白字符分割

		// 拼接
		"join": strings.Join, // 用分隔符将切片拼接为字符串

		// 大小写转换
		"toLower":     strings.ToLower,     // 转为全小写
		"toUpper":     strings.ToUpper,     // 转为全大写
		"toTitle":     strings.ToTitle,     // 转为标题（Title Case）
		"toValidUTF8": strings.ToValidUTF8, // 将无效 UTF-8 替换为指定字符串

		// 裁剪
		"trim":       strings.Trim,       // 去除两端指定字符
		"trimLeft":   strings.TrimLeft,   // 去除左侧指定字符
		"trimRight":  strings.TrimRight,  // 去除右侧指定字符
		"trimSpace":  strings.TrimSpace,  // 去除两端空白字符
		"trimPrefix": strings.TrimPrefix, // 去除前缀（不匹配则返回原字符串）
		"trimSuffix": strings.TrimSuffix, // 去除后缀（不匹配则返回原字符串）

		// 比较
		"compare":   strings.Compare,   // 字典序比较（返回 -1/0/+1）
		"equalFold": strings.EqualFold, // 忽略大小写比较

		// 构造器
		"newReader":   strings.NewReader,   // 将字符串包装为 io.Reader
		"newReplacer": strings.NewReplacer, // 创建多对替换器（old1, new1, old2, new2, ...）
	}

	return jsvm.ModuleExports{
		Name:    "strings",
		Default: vals,
	}
}
