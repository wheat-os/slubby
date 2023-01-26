package download

import (
	"github.com/wheat-os/slubby/v2/stream"
	"github.com/wheat-os/slubby/v2/stream/buffer"
)

type Option func(opt *SlubbyDownload)

type retryContentKey string

func withRetry(retryNum int, to stream.Cover) Option {
	const key retryContentKey = "slubby.stream.httpdownload.retry"
	return func(opt *SlubbyDownload) {
		opt.isRetryFunc = func(req stream.Stream, resp stream.Stream) (stream.Cover, bool) {
			if retryNum <= 0 {
				return stream.UnknownCover, false
			}

			if req == nil {
				return stream.UnknownCover, false
			}
			numV := req.GetMeta(key)

			switch num := numV.(type) {
			case int:
				req.SetMeta(key, num+1)
				return to, retryNum-num >= 0
			case nil:
				req.SetMeta(key, 1)
				return to, retryNum > 0
			}

			return stream.UnknownCover, false
		}

	}
}

// WithDirectRetry 使用立即重试方案
func WithDirectRetry(retryNum int) Option {
	return withRetry(retryNum, stream.DownloadCover)
}

// WithBackSchedulerRetry 回调度器重试方案
func WithBackSchedulerRetry(retryNum int) Option {
	return withRetry(retryNum, stream.SchedulerCover)
}

// WithRetryFunc 自定义重试检查
func WithRetryFunc(fn func(req stream.Stream, resp stream.Stream) (stream.Cover, bool)) Option {
	return func(opt *SlubbyDownload) {
		opt.isRetryFunc = fn
	}
}

// WithTransport 自定义 transport
func WithTransport(transport RoundTripper) Option {
	return func(opt *SlubbyDownload) {
		opt.roundTripper = transport
	}
}

// WithForwardCover 定义下载成功转发
func WithForwardCover(cover stream.Cover) Option {
	return func(opt *SlubbyDownload) {
		opt.forwardCover = cover
	}
}

// WithDownloadBuffer 定义下载器缓冲区
func WithDownloadBuffer(buffer buffer.StreamBuffer) Option {
	return func(opt *SlubbyDownload) {
		opt.buffer = buffer
	}
}

// WithDownloadProcess 设置下载器进程数量
func WithDownloadProcess(process int) Option {
	return func(opt *SlubbyDownload) {
		opt.process = process
	}
}
