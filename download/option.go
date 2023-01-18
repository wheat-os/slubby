package download

import "github.com/wheat-os/slubby/v2/stream"

type OptFunc func(opt *SlubbyComponent)

type RoundTripper func(inStream stream.Stream) (stream.Stream, error)

type retryContentKey string

func withRetry(retryNum int, to stream.Cover) OptFunc {
	const key retryContentKey = "slubby.stream.httpdownload.retry"
	return func(opt *SlubbyComponent) {
		opt.isRetryFunc = func(req stream.Stream, resp stream.Stream) (stream.Cover, bool) {
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
func WithDirectRetry(retryNum int) OptFunc {
	return withRetry(retryNum, stream.DownloadCover)
}

// WithBackSchedulerRetry 回调度器重试
func WithBackSchedulerRetry(retryNum int) OptFunc {
	return withRetry(retryNum, stream.SchedulerCover)
}

// WithRetryFunc 自定义重试检查
func WithRetryFunc(fn func(req stream.Stream, resp stream.Stream) (stream.Cover, bool)) OptFunc {
	return func(opt *SlubbyComponent) {
		opt.isRetryFunc = fn
	}
}

// WithTransport 自定义 transport
func WithTransport(transport RoundTripper) OptFunc {
	return func(opt *SlubbyComponent) {
		opt.roundTripper = transport
	}
}

// WithForwardCover 定义下载成功转发
func WithForwardCover(cover stream.Cover) OptFunc {
	return func(opt *SlubbyComponent) {
		opt.forwardCover = cover
	}
}
