package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

// HttpRequest 函数用于发送HTTP请求
func HttpRequest(
	uslStr string,
	method string,
	headers map[string]string,
	params map[string]string, // 查询参数（如 ?key=value&key2=value2）
	data any) (*http.Response, error) { // 请求体的内容（如果有的话）
	u, err := url.Parse(uslStr)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode() // 更新URL的查询部分
	// 将请求体数据（如果有）编码成JSON格式
	buf := new(bytes.Buffer) // 创建一个缓冲区用于存储请求体
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b) // 将编码后的字节数组转换为缓冲区
	}

	// 创建HTTP请求对象
	req, err := http.NewRequest(method, u.String(), buf) // 使用指定的URL和方法创建请求
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// 发送HTTP请求并获取响应
	resp, err := http.DefaultClient.Do(req) // 使用默认的HTTP客户端发送请求
	if err != nil {
		return nil, err
	}

	// 返回响应对象
	return resp, nil // 返回响应，调用者可根据需要处理响应数据
}
