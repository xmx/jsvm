package jshttp

import (
	"encoding/json"
	"io"
	"net/http"
)

// response 封装标准 net/http.Response，通过 goja tagMapper 自动将方法
// 名映射为小写（JSON → json，Text → text）供 JS 调用。
type response struct {
	*http.Response
}

// JSON 将响应体解析为任意 JSON 值并返回。
// 解码后 Body 已被消费，不可再次读取。
func (r *response) JSON() (any, error) {
	var ret any
	dec := json.NewDecoder(r.Response.Body)
	if err := dec.Decode(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// Text 将响应体全部读取为字符串返回。
// 读取后 Body 已被消费，不可再次读取。
func (r *response) Text() (string, error) {
	bs, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}
