package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/wheat-os/slubby/v2/stream"
)

type StreamResponse struct {
	stream.TargetCover
	Response *http.Response
	stream.Context
}

// Read 读取 response body 消息
func (s *StreamResponse) Read(p []byte) (n int, err error) {
	if s.Response == nil || s.Response.Body == nil {
		return 0, ErrNotStreamBody
	}
	return s.Response.Body.Read(p)
}

// Write 写入 response body 消息
func (s *StreamResponse) Write(p []byte) (n int, err error) {
	if s.Response == nil || s.Response.Body == nil {
		return 0, ErrNotStreamBody
	}
	// 可用的 http body
	if wr, ok := s.Response.Body.(io.Writer); ok {
		return wr.Write(p)
	}

	wr := bytes.NewBuffer(nil)
	if _, err := io.Copy(wr, s.Response.Body); err != nil {
		return 0, err
	}
	_ = s.Response.Body.Close()
	s.Response.Body = io.NopCloser(wr)
	return wr.Write(p)
}

// Close 关闭 response
func (s *StreamResponse) Close() error {
	if s.Response != nil || s.Response.Body == nil {
		return nil
	}
	return s.Response.Body.Close()
}

// ReplaceCtx 替换或者获取上下文
func (s *StreamResponse) ReplaceCtx(ctx stream.Context) stream.Context {
	vtr := s.Context
	if ctx != nil {
		s.Context = vtr
	}
	return vtr
}

// NewHttpResponse 新创建 http response
func NewHttpResponse(resp *http.Response, meta stream.Context) (*StreamResponse, error) {
	if meta == nil {
		meta = &stream.Meta{}
	}
	return &StreamResponse{Response: resp, Context: meta}, nil
}
