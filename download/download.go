package download

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/wheat-os/slubby/v2/engine"
	"github.com/wheat-os/slubby/v2/stream"
	"github.com/wheat-os/slubby/v2/stream/buffer"
)

var (
	ErrRequestBufferIsClose = errors.New("request buffer is close")
)

type RoundTripper func(inStream stream.Stream) (stream.Stream, error)

type SlubbyDownload struct {
	roundTripper RoundTripper
	// 返回 cover 以及 bool, 标识是否重试, 以及转 component 处理,
	// stream.Cover 为 stream.DownloadCover 时 直接重试
	isRetryFunc  func(req stream.Stream, resp stream.Stream) (stream.Cover, bool)
	forwardCover stream.Cover // 默认转发组件

	buffer    buffer.StreamBuffer             // 下载器缓冲区
	recv      chan stream.Stream              // 通信队列
	rateLimit func(ctx context.Context) error // 下载器限制器

	process int // 下载器线程数

	isClose bool
	cancel  func()
}

func (h *SlubbyDownload) pushRequest(ctx context.Context, data stream.Stream) error {
	if h.isClose {
		return ErrRequestBufferIsClose
	}

	return h.buffer.PutStream(ctx, data)
}

func (h *SlubbyDownload) pullRequest(ctx context.Context) (stream.Stream, error) {
	return h.buffer.GetStream(ctx)
}

func (h *SlubbyDownload) Size() int {
	return h.buffer.Len()
}

// downTripper 实现请求的首发方法
func (h *SlubbyDownload) downTripper(req stream.Stream) (stream.Stream, stream.Cover, error) {
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

// do 工作流
func (h *SlubbyDownload) do(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			h.recv <- stream.FromError(stream.DownloadCover, fmt.Errorf("%v", err))
		}
	}()

	if h.Size() == 0 {
		return
	}

	req, err := h.pullRequest(ctx)
	switch err {
	case nil:
		// 忽略读取错误
	case buffer.ErrStreamBufferIsEmpty:
		return
	default:
		h.recv <- stream.FromError(stream.DownloadCover, err)
		return
	}

	// 执行限流器
	if err := h.rateLimit(ctx); err != nil {
		h.recv <- stream.Error(err)
		req.SetForm(stream.DownloadCover)
		req.SetTo(stream.SchedulerCover)
		h.recv <- req
		return
	}

	// 非下载器处理流处理
	if req.To() != stream.DownloadCover {
		req.SetForm(stream.DownloadCover)
		h.recv <- req
		return
	}

	// 执行 downTripper
	fromStream, cover, err := h.downTripper(req)
	if err != nil {
		h.recv <- stream.FromError(stream.DownloadCover, err)
		return
	}

	// 发送 cover
	fromStream.SetForm(stream.DownloadCover)
	fromStream.SetTo(cover)
	h.recv <- fromStream
}

// working 工作方法
func (h *SlubbyDownload) working(ctx context.Context) {
	for {
		// 检查关闭
		select {
		case <-ctx.Done():
			return
		default:
		}

		h.do(ctx)
	}
}

func (h *SlubbyDownload) run() func() {
	// 运行 slubby download
	group := sync.WaitGroup{}
	cancelFunc := make([]func(), 0, h.process)

	for i := 0; i < h.process; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancelFunc = append(cancelFunc, cancel)

		// 执行 download 进程
		group.Add(1)
		go func(c context.Context) {
			h.working(c)
			group.Done()
		}(ctx)
	}

	// 生成关闭方法
	return func() {
		for _, cancel := range cancelFunc {
			cancel()
		}
		group.Wait()
	}
}

// Streaming 接受下载流
func (h *SlubbyDownload) Streaming(data stream.Stream) error {
	return h.pushRequest(context.TODO(), data)
}

// BackStream 获取响应推流器
func (h *SlubbyDownload) BackStream() <-chan stream.Stream {
	return h.recv
}

// Close 关闭下载器
func (h *SlubbyDownload) Close() error {
	h.isClose = true
	h.cancel()
	return nil
}

// NewSlubbyDownload 创建一个默认的 slubby 下载器
func NewSlubbyDownload(opts ...Option) engine.SendAndReceiveComponent {
	var defaultProcessNum = runtime.NumCPU()

	component := &SlubbyDownload{recv: make(chan stream.Stream)}

	initOption := []Option{
		WithHttpClientTransport(&http.Client{}),                     // 默认使用默认的 HTTP 下载器
		WithDownloadProcess(defaultProcessNum),                      // 默认使用 CPU 数量的进程数
		WithDirectRetry(0),                                          // 默认不执行重试
		WithForwardCover(stream.SpiderCover),                        // 默认下载成功后转发 spider 处理
		WithDownloadBuffer(buffer.NewChanBuffer(defaultProcessNum)), // 默认使用 chan 缓存队列
		WithTokenBucketLimit(-1, 1),                                 // 默认不启用限制器
	}

	initOption = append(initOption, opts...)
	for _, optFunc := range initOption {
		optFunc(component)
	}

	workingCancel := component.run()

	// 写入关闭函数
	component.cancel = func() {
		workingCancel()
	}

	return component
}
