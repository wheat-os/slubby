package scheduler

import (
	"errors"
	"github.com/wheat-os/slubby/v2/stream"
	"github.com/wheat-os/slubby/v2/stream/buffer"
)

var (
	ErrSchedulerIsClose = errors.New("request buffer is close")
)

type SlubbyScheduler struct {
	buffer     buffer.StreamBuffer          // 调度器缓冲区
	passFilter func(stm stream.Stream) bool // 过滤器机制，是否通过过滤器

	recv chan stream.Stream

	isClose bool
}

func (s *SlubbyScheduler) pushStream(stm stream.Stream) error {
	if s.isClose {
		return ErrSchedulerIsClose
	}

	// 检查过滤器情况
	if !s.passFilter(stm) {
		return nil
	}

	return s.buffer.PutStream(stm)
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
	return nil
}
