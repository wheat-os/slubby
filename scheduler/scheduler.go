package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/wheat-os/slubby/v2/engine"
	"github.com/wheat-os/slubby/v2/stream"
	"github.com/wheat-os/slubby/v2/stream/buffer"
	"io"
	"sync"
)

var (
	ErrSchedulerIsClose  = errors.New("scheduler is close")
	ErrSchedulerIsFinish = errors.New("scheduler is finish")
)

// StreamFilter 流过滤器
type StreamFilter interface {
	PassFilter(b []byte) bool
	io.Closer
}

type SlubbyScheduler struct {
	buffer       buffer.StreamBuffer // 调度器缓冲区
	filter       StreamFilter        // 过滤器机制
	enc          stream.Encoder      // 流编码器
	forwardCover stream.Cover        // 默认转发组件

	recv    chan stream.Stream
	process int // 执行线程数

	isClose  bool
	isFinish bool
	cancel   func()
}

func (s *SlubbyScheduler) pushStream(stm stream.Stream) error {
	// 调度器持久化完成
	if s.isFinish {
		return ErrSchedulerIsClose
	}

	buf, err := s.enc.StreamEncode(stm)
	if err != nil {
		return err
	}

	// 不通过过滤器
	if !s.filter.PassFilter(buf) {
		return nil
	}

	return s.buffer.PutStream(context.TODO(), stm)
}

func (s *SlubbyScheduler) pullStream(ctx context.Context) (stream.Stream, error) {
	if s.isClose {
		return nil, ErrSchedulerIsClose
	}

	return s.buffer.GetStream(ctx)
}

func (s *SlubbyScheduler) do(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			s.recv <- stream.FromError(stream.SchedulerCover, fmt.Errorf("%v", err))
		}
	}()

	if s.buffer.Len() == 0 {
		return
	}

	stm, err := s.pullStream(ctx)
	switch err {
	case nil:
		// 忽略空读错误
	case buffer.ErrStreamBufferIsEmpty, buffer.ErrStreamContentIsCancel:
		return
	default:
		s.recv <- stream.FromError(stream.SchedulerCover, err)
		return
	}

	// 默认推送
	stm.SetForm(stream.SchedulerCover)
	stm.SetForm(s.forwardCover)
	s.recv <- stm
}

func (s *SlubbyScheduler) working(ctx context.Context) {
	for {
		// 检查关闭
		select {
		case <-ctx.Done():
			return
		default:
		}
		s.do(ctx)
	}
}

// run 创建工作流
func (s *SlubbyScheduler) run() func() {
	group := sync.WaitGroup{}
	cancelFunc := make([]func(), 0, 1)

	for i := 0; i < s.process; i++ {
		ctx, cancel := context.WithCancel(context.TODO())

		// 执行 download 进程
		group.Add(1)
		go func(c context.Context) {
			s.working(c)
			group.Done()
		}(ctx)
		cancelFunc = append(cancelFunc, cancel)
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
func (s *SlubbyScheduler) Streaming(data stream.Stream) error {
	return s.pushStream(data)
}

// BackStream 获取响应推流器
func (s *SlubbyScheduler) BackStream() <-chan stream.Stream {
	return s.recv
}

// Close 关闭下载器
func (s *SlubbyScheduler) Close() error {
	s.isClose = true
	s.cancel()

	// 关闭对外推送管道
	close(s.recv)

	return nil
}

// Finish 结束时持久化调用
func (s *SlubbyScheduler) Finish() error {
	s.isFinish = true

	// 执行持久化
	return s.filter.Close()
}

// NewSlubbyScheduler 创建一个 slubby 调度器
func NewSlubbyScheduler(opts ...Option) engine.SchedulerComponent {
	component := &SlubbyScheduler{recv: make(chan stream.Stream)}

	initOption := []Option{
		WithUselessFilter(),                              // 默认不使用过滤器
		WithProcess(1),                                   // 默认使用单个进程
		WithForwardCover(stream.DownloadCover),           // 定义默认流出口
		WithStreamEncoder(stream.NewBackGroundEncoder()), // 默认编码器
		WithSchedulerBuffer(buffer.NewListBuffer(10000)), // 默认储存 1w 个调度单位
	}

	initOption = append(initOption, opts...)
	for _, optFunc := range initOption {
		optFunc(component)
	}

	// 写入关闭函数
	workingCancel := component.run()
	component.cancel = func() {
		workingCancel()
	}

	return component
}
