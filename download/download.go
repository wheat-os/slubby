package download

import (
	"net/http"

	"gitee.com/wheat-os/slubby/download/middle"
	"gitee.com/wheat-os/slubby/stream"
	"github.com/panjf2000/ants/v2"
)

type Download interface {
	middle.Middleware
	Do(req *stream.HttpRequest) (*stream.HttpResponse, error)
	Activate() bool

	Close() error
}

type shortDownload struct {
	middle.Middleware

	// 配置文件
	opt *option
}

func (s *shortDownload) poll() *ants.Pool {
	return s.opt.poll
}

func (s *shortDownload) client() *http.Client {
	return s.opt.client
}

func (s *shortDownload) Do(req *stream.HttpRequest) (*stream.HttpResponse, error) {

	var (
		resp *http.Response
		err  error
	)

	m := middle.MC()
	if req, err = s.BeforeDownload(m, req); err != nil {
		return nil, err
	}

	ch := make(chan struct{})
	s.poll().Submit(func() {
		resp, err = s.client().Do(req.Request)
		ch <- struct{}{}
	})
	<-ch

	if err != nil {
		s.ProcessErr(m, req, err)
	}

	eResp := new(stream.HttpResponse)
	eResp.WithHttpAndRequestStream(req, resp)

	if eResp, err = s.AfterDownload(m, req, eResp); err != nil {
		return nil, err
	}

	return eResp, nil
}

func (s *shortDownload) BeforeDownload(
	m *middle.M,
	req *stream.HttpRequest,
) (*stream.HttpRequest, error) {
	if s.Middleware == nil {
		return req, nil
	}
	m.Next()
	return s.Middleware.BeforeDownload(m, req)
}

func (s *shortDownload) AfterDownload(
	m *middle.M,
	req *stream.HttpRequest,
	resp *stream.HttpResponse,
) (*stream.HttpResponse, error) {
	if s.Middleware == nil {
		return resp, nil
	}
	m.Next()
	return s.Middleware.AfterDownload(m, req, resp)
}

func (s *shortDownload) ProcessErr(
	m *middle.M,
	req *stream.HttpRequest,
	err error,
) {
	if s.Middleware != nil {
		m.Next()
		s.Middleware.ProcessErr(m, req, err)
	}
}

func (s *shortDownload) Activate() bool {
	return s.poll().Running() == 0
}

func (s *shortDownload) Close() error {
	s.poll().Release()
	return nil
}

func DefaultDownload() Download {
	return &shortDownload{
		opt: loadOption(),
	}
}

func NewShortDownload(mid middle.Middleware, opts ...optionFunc) Download {
	return &shortDownload{
		opt:        loadOption(opts...),
		Middleware: mid,
	}
}
