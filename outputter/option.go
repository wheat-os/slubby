package outputter

import (
	"runtime"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/wheat-os/slubby/outputter/pipline"
	"github.com/wheat-os/slubby/pkg/tools"
)

type option struct {
	threadCount int
	pip         pipline.Pipline
	poll        *ants.Pool
	once        sync.Once
}

// creat and init Poll
func (o *option) Poll() *ants.Pool {
	o.once.Do(func() {
		poll, err := ants.NewPool(
			o.threadCount,
			ants.WithPanicHandler(tools.AntsWlogHandlePanic),
		)
		if err != nil {
			panic(err)
		}
		o.poll = poll
	})

	return o.poll
}

type optionFunc func(o *option)

// do pipline
func WithPipline(pip pipline.Pipline) optionFunc {
	return func(o *option) {
		o.pip = pip
	}
}

func WithGroupPipline(pip ...pipline.Pipline) optionFunc {
	return func(o *option) {
		o.pip = pipline.GroupPipline(pip...)
	}
}

func WithThreadCount(num int) optionFunc {
	return func(o *option) {
		o.threadCount = num
	}
}

func loadOPtion(opts ...optionFunc) *option {
	ops := &option{
		threadCount: runtime.NumCPU(),
	}

	ops.Poll()

	for _, opt := range opts {
		opt(ops)
	}

	return ops
}
