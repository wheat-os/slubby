package engine

import (
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wheat-os/slubby/download"
	"github.com/wheat-os/slubby/outputter"
	"github.com/wheat-os/slubby/pkg/sundry"
	"github.com/wheat-os/slubby/scheduler"
	"github.com/wheat-os/slubby/spider"
	"github.com/wheat-os/wlog"
)

type option struct {
	threadCount int

	// engine 检查时间
	checkTime time.Duration
	poll      *ants.Pool
	once      sync.Once

	// 基础组件
	eDownload  download.Download
	eScheduler scheduler.Scheduler
	eOutPutter outputter.Outputter
	eSpider    *spider.Manage
}

func (o *option) Poll() *ants.Pool {
	o.once.Do(func() {
		poll, err := ants.NewPool(
			o.threadCount,
			ants.WithPanicHandler(sundry.AntsWlogHandlePanic),
		)
		if err != nil {
			wlog.Panic(err)
		}

		o.poll = poll
	})

	return o.poll
}

type optionFunc func(o *option)

// 默认
func loadOption(opts ...optionFunc) *option {
	ops := &option{
		threadCount: 64,
		checkTime:   time.Second,

		eDownload:  download.ShortDownload(),
		eScheduler: scheduler.ShortScheduler(),
		eOutPutter: outputter.ShortOutputter(),
		eSpider:    spider.SpiderManage(),
	}

	ops.Poll()

	for _, opt := range opts {
		opt(ops)
	}

	return ops
}

func WithDownload(down download.Download) optionFunc {
	return func(o *option) {
		o.eDownload = down
	}
}

func WithScheduler(scd scheduler.Scheduler) optionFunc {
	return func(o *option) {
		o.eScheduler = scd
	}
}

func WithTreadCount(num int) optionFunc {
	return func(o *option) {
		o.threadCount = num
	}
}

func WithOutputter(out outputter.Outputter) optionFunc {
	return func(o *option) {
		o.eOutPutter = out
	}
}

func WithCheckTime(t time.Duration) optionFunc {
	return func(o *option) {
		o.checkTime = t
	}
}
