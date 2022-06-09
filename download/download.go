package download

import (
	"net/http"

	"github.com/panjf2000/ants/v2"
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/stream"
)

type Download interface {
	middle.Middleware
	Do(req *stream.HttpRequest) (*stream.HttpResponse, error)
	Activate() bool

	Close() error
}

type shortDownload struct {
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
	if req, err = s.BeforeDownload(m, req); err != nil || req == nil {
		return nil, err
	}

	ch := make(chan struct{})
	s.poll().Submit(func() {

		// 获取 request 令牌
		if s.opt.limiter != nil {
			s.opt.limiter.Allow(req)
		}

		for req.Retry > 0 {
			if resp, err = s.client().Do(req.Request); err == nil {
				break
			}
			req.Retry -= 1
		}

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
	if s.opt.Middleware == nil {
		return req, nil
	}
	m.Next()
	return s.opt.Middleware.BeforeDownload(m, req)
}

func (s *shortDownload) AfterDownload(
	m *middle.M,
	req *stream.HttpRequest,
	resp *stream.HttpResponse,
) (*stream.HttpResponse, error) {
	if s.opt.Middleware == nil {
		return resp, nil
	}
	m.Next()
	return s.opt.Middleware.AfterDownload(m, req, resp)
}

func (s *shortDownload) ProcessErr(
	m *middle.M,
	req *stream.HttpRequest,
	err error,
) {
	if s.opt.Middleware != nil {
		m.Next()
		s.opt.Middleware.ProcessErr(m, req, err)
	}
}

func (s *shortDownload) Activate() bool {
	return s.poll().Running() != 0
}

func (s *shortDownload) Close() error {
	s.poll().Release()
	return nil
}

func ShortDownload(opt ...optionFunc) Download {
	return &shortDownload{
		opt: loadOption(opt...),
	}
}
