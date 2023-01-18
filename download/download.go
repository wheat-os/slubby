package download

import (
	"context"
	"errors"
	"github.com/wheat-os/slubby/v2/stream"
	"github.com/wheat-os/slubby/v2/stream/buffer"
)

var (
	ErrRequestBufferIsClose = errors.New("request buffer is close")
)

type SlubbyComponent struct {
	roundTripper RoundTripper
	// 返回 cover 以及 bool, 标识是否重试, 以及转 component 处理,
	// stream.Cover 为 stream.DownloadCover 时 直接重试
	isRetryFunc  func(req stream.Stream, resp stream.Stream) (stream.Cover, bool)
	forwardCover stream.Cover // 默认转发组件

	buffer buffer.StreamBuffer

	recv    chan stream.Stream
	isClose bool
}

func (h *SlubbyComponent) pushRequest(data stream.Stream) error {
	if h.isClose {
		return ErrRequestBufferIsClose
	}

	return h.buffer.PutStream(data)
}

func (h *SlubbyComponent) pullRequest() (stream.Stream, error) {
	return h.buffer.GetStream()
}

func (h *SlubbyComponent) Size() int {
	return h.buffer.Len()
}

// downTripper 实现请求的首发方法
func (h *SlubbyComponent) downTripper(req stream.Stream) (stream.Stream, stream.Cover, error) {
	for {
		resp, err := h.roundTripper(req)
		if err == nil {
			return resp, h.forwardCover, nil
		}

		cover, b := h.isRetryFunc(req, resp)

		// 超过重试次数返回
		if !b {
			return nil, cover, err
		}

		switch cover {
		// 默认转下载器继续处理
		case stream.DownloadCover, stream.UnknownCover:
			continue
		default:
			return req, cover, nil
		}
	}
}

// working 工作流
func (h *SlubbyComponent) working(ctx context.Context) {
	for {
		// 检查关闭
		select {
		case <-ctx.Done():
			return
		default:
		}

		if h.Size() == 0 {
			continue
		}

		req, err := h.pullRequest()
		switch err {
		case nil:
			// 忽略读取错误
		case buffer.ErrStreamBufferIsEmpty:
			continue
		default:
			h.recv <- stream.FromError(stream.DownloadCover, err)
			continue
		}

		// 非下载器处理流处理
		if req.To() != stream.DownloadCover {
			req.SetForm(stream.DownloadCover)
			h.recv <- req
			continue
		}

		// 执行 downTripper
		fromStream, cover, err := h.downTripper(req)
		if err != nil {
			h.recv <- stream.FromError(stream.DownloadCover, err)
			continue
		}

		// 发送 cover
		fromStream.SetForm(stream.DownloadCover)
		fromStream.SetTo(cover)
		h.recv <- fromStream
	}
}

// Streaming 接受下载流
func (h *SlubbyComponent) Streaming(data stream.Stream) error {
	return h.pushRequest(data)
}

// BackStream 获取响应推流器
func (h *SlubbyComponent) BackStream() <-chan stream.Stream {
	return h.recv
}

// Close 关闭下载器
func (h *SlubbyComponent) Close() error {
	//TODO implement me
	panic("implement me")
}
