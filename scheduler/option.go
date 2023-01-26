package scheduler

import (
	"github.com/wheat-os/slubby/v2/stream"
	"github.com/wheat-os/slubby/v2/stream/buffer"
)

type Option func(s *SlubbyScheduler)

// WithProcess 指定执行进程数
func WithProcess(process int) Option {
	return func(s *SlubbyScheduler) {
		s.process = process
	}
}

// WithSchedulerBuffer 指定调度器缓冲区
func WithSchedulerBuffer(buffer buffer.StreamBuffer) Option {
	return func(s *SlubbyScheduler) {
		s.buffer = buffer
	}
}

// WithForwardCover 定义下载成功转发
func WithForwardCover(cover stream.Cover) Option {
	return func(opt *SlubbyScheduler) {
		opt.forwardCover = cover
	}
}
