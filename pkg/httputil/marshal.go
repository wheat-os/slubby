package httputil

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httputil"
	"strings"
)

// DumpRequest 编码 http 请求
func DumpRequest(req *http.Request, body bool) ([]byte, error) {
	return httputil.DumpRequest(req, body)
}

// LoadRequest 解编码 http 请求
func LoadRequest(buf []byte) (*http.Request, error) {
	record := strings.Split(string(buf), "\r\n")
	// URI 解析
	var (
		url        string
		method     string
		body       = bytes.NewBuffer(nil)
		parseForce int
		headers    = make(map[string]string)
	)
	for _, force := range record {
		switch parseForce {
		// 解析头部消息
		case 0:
			meta := strings.Split(force, " ")
			if len(meta) != 3 {
				return nil, errors.New("failed to parse the request line")
			}
			method = meta[0]
			url = meta[1]
			parseForce++

		// 解析 headers
		case 1:
			if len(force) == 0 {
				parseForce++
				continue
			}

			parm := strings.Split(force, ":")
			if len(parm) != 2 {
				return nil, errors.New("failed to parse the request line")
			}
			headers[parm[0]] = parm[1]

		// 分配 body 空间, 初始化 body
		case 2:
			body.Grow(len(record[len(record)-1]))
			body.WriteString(force)
			parseForce++

		case 3:
			body.WriteString("\r\n")
			body.WriteString(force)
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// write headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}
