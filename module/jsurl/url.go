package jsurl

import (
	"net/url"

	"github.com/dop251/goja"
	"github.com/xmx/jsvm"
)

// urlModule 是 net/url 模块的实现，将 Go url 包的解析与构建功能导出到 JS。
type urlModule struct{}

// New 创建 net/url 模块实例。
func New() jsvm.ModuleExporter {
	return &urlModule{}
}

// ModuleExports 注册 net/url 模块，导出：
//   - parse/parseRequestURI：解析 URL 字符串
//   - parseQuery：解析查询字符串为 Values 对象
//   - queryEscape/queryUnescape：URL 编码与解码
//   - user/userPassword：构建用户认证信息
//   - Values：Values 对象的构造函数
func (m *urlModule) ModuleExports(*jsvm.VM) jsvm.ModuleExports {
	defaults := map[string]any{
		"parse":           url.Parse,           // 解析 URL 字符串，返回 URL 对象和错误
		"parseRequestURI": url.ParseRequestURI, // 解析请求 URI（不允许相对 URL）
		"parseQuery":      url.ParseQuery,      // 解析查询字符串为 Values
		"queryEscape":     url.QueryEscape,     // 对字符串进行 URL 百分号编码
		"queryUnescape":   url.QueryUnescape,   // 对 URL 编码字符串进行解码
		"user":            url.User,            // 创建仅含用户名的 Userinfo
		"userPassword":    url.UserPassword,    // 创建含用户名和密码的 Userinfo
		"Values":          m.newValues,         // Values 构造函数
	}

	return jsvm.ModuleExports{
		Name:    "net/url",
		Default: defaults,
	}
}

// newValues 是 url.Values 的构造函数，供 JS 使用 `new url.Values()` 语法创建。
// Values 本质是 map[string][]string，用于构建或解析 URL 查询参数。
func (m *urlModule) newValues(_ goja.ConstructorCall, rt *goja.Runtime) *goja.Object {
	v := make(url.Values)
	val := rt.ToValue(v)

	return val.ToObject(rt)
}
