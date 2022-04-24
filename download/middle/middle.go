package middle

import "gitee.com/wheat-os/slubby/stream"

const (
	nextMiddleM = iota
	abortMiddleM
)

// 中间件流程控制器
type M uint8

func (m *M) Abort() {
	*m = abortMiddleM
}

func (m *M) Next() {
	*m = nextMiddleM
}

// Relay 检查状态， 向下
func (m *M) Relay() bool {
	return *m == nextMiddleM
}

func MC() *M {
	return new(M)
}

type Middleware interface {
	BeforeDownload(m *M, req *stream.HttpRequest) (*stream.HttpRequest, error)

	AfterDownload(m *M, req *stream.HttpRequest, resp *stream.HttpResponse) (*stream.HttpResponse, error)

	ProcessErr(m *M, req *stream.HttpRequest, err error)
}

type MiddlewareGroup []Middleware

func NewMiddleGroup(m ...Middleware) MiddlewareGroup { return m }

func (mid MiddlewareGroup) BeforeDownload(m *M, req *stream.HttpRequest) (*stream.HttpRequest, error) {

	var err error
	for _, handle := range mid {
		// 检查是否向下中继
		if !m.Relay() {
			break
		}
		req, err = handle.BeforeDownload(m, req)
	}

	return req, err
}

func (mid MiddlewareGroup) AfterDownload(m *M, req *stream.HttpRequest, resp *stream.HttpResponse) (*stream.HttpResponse, error) {
	var err error
	for _, handle := range mid {
		// 检查是否向下中继
		if !m.Relay() {
			break
		}
		resp, err = handle.AfterDownload(m, req, resp)
	}

	return resp, err
}

func (mid MiddlewareGroup) ProcessErr(m *M, req *stream.HttpRequest, resp *stream.HttpResponse, err error) {

	for _, handle := range mid {
		// 检查是否向下中继
		if !m.Relay() {
			break
		}
		handle.ProcessErr(m, req, err)
	}

}
