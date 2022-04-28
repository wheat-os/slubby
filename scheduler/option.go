package scheduler

import (
	"gitee.com/wheat-os/slubby/scheduler/buffer"
	"gitee.com/wheat-os/slubby/scheduler/filter"
)

type optionFunc = func(o *option)

type option struct {
	// 过滤器
	filter filter.Filter

	// 缓冲区
	buffer buffer.Buffer
}

func WithFilter(filter filter.Filter) optionFunc {
	return func(o *option) {
		o.filter = filter
	}
}

func WithBuffer(buf buffer.Buffer) optionFunc {
	return func(o *option) {
		o.buffer = buf
	}
}

// load option
func loadOption(opts ...optionFunc) *option {
	opt := &option{
		buffer: buffer.ShortQueue(),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	return opt
}
