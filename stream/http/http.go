package http

import (
	"bytes"
	"errors"
	"github.com/wheat-os/slubby/v2/stream"
	"io"
	"net/http"
	"net/http/httputil"
)

var (
	ErrNotStreamBody = errors.New("this http request stream not have body")
)

type StreamRequest struct {
	stream.TargetCover
	Request *http.Request
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

// MarshalBinary 编码 http 请求使用字符串格式
func (s *StreamRequest) MarshalBinary() (data []byte, err error) {
	return httputil.DumpRequest(s.Request, true)
}

type StreamResponse struct {
	stream.TargetCover
	Response *http.Response
}

// Read 读取 response body 消息
func (s *StreamResponse) Read(p []byte) (n int, err error) {
	if s.Response != nil || s.Response.Body == nil {
		return 0, ErrNotStreamBody
	}
	return s.Response.Body.Read(p)
}

// Write 写入 response body 消息
func (s *StreamResponse) Write(p []byte) (n int, err error) {
	if s.Response != nil || s.Response.Body == nil {
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
	return s.Response.Body.Close()
}

// MarshalBinary 编码 response
func (s *StreamResponse) MarshalBinary() (data []byte, err error) {
	return httputil.DumpResponse(s.Response, true)
}
