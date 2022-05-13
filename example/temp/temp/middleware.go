package temp

import (
	"gitee.com/wheat-os/slubby/download/middle"
	"gitee.com/wheat-os/slubby/stream"
)

type SimpleMiddle struct{}

func (s *SimpleMiddle) BeforeDownload(m *middle.M, req *stream.HttpRequest) (*stream.HttpRequest, error) {
	return req, nil
}

func (s *SimpleMiddle) AfterDownload(
	m *middle.M,
	req *stream.HttpRequest,
	resp *stream.HttpResponse,
) (*stream.HttpResponse, error) {
	return resp, nil
}

func (s *SimpleMiddle) ProcessErr(m *middle.M, req *stream.HttpRequest, err error) {}
