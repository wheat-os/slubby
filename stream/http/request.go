package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/wheat-os/slubby/v2/stream"
)

type StreamRequest struct {
	stream.TargetCover
	Request *http.Request
	stream.Context
}

// Read 读取 http request body 的消息
func (s *StreamRequest) Read(p []byte) (n int, err error) {
	if s.Request == nil || s.Request.Body == nil {
		return 0, ErrNotStreamBody
	}

	return s.Request.Body.Read(p)
}

// Write 写入 http body
func (s *StreamRequest) Write(p []byte) (n int, err error) {
	if s.Request == nil || s.Request.Body == nil {
		return 0, ErrNotStreamBody
	}

	// 可用的 http body
	if wr, ok := s.Request.Body.(io.Writer); ok {
		return wr.Write(p)
	}

	wr := bytes.NewBuffer(nil)
	if _, err := io.Copy(wr, s.Request.Body); err != nil {
		return 0, err
	}
	s.Request.Body = io.NopCloser(wr)

	return wr.Write(p)
}

// Close 关闭 http 请求
func (s *StreamRequest) Close() error {
	if s.Request == nil || s.Request.Body == nil {
		return ErrNotStreamBody
	}

	return s.Request.Body.Close()
}

// ReplaceCtx 替换或者获取上下文
func (s *StreamRequest) ReplaceCtx(ctx stream.Context) stream.Context {
	vtr := s.Context
	if ctx != nil {
		s.Context = vtr
	}
	return vtr
}

// NewHttpRequest 创建 http request
func NewHttpRequest(req *http.Request) (*StreamRequest, error) {
	return &StreamRequest{Request: req, Context: &stream.Meta{}}, nil
}
