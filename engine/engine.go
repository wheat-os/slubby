package engine

import (
	"context"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wheat-os/slubby/spider"
	"github.com/wheat-os/slubby/stream"
	"github.com/wheat-os/wlog"
)

type Engine interface {
	Start(ctx context.Context)
	Close()

	Register(sp spider.Spider) error
}

type shortEngine struct {
	opt     *option
	spiders []spider.Spider
}

func (e *shortEngine) poll() *ants.Pool {
	return e.opt.poll
}

func (e *shortEngine) canClose() bool {
	// engine 是否可以关闭
	if e.poll().Running() != 0 {
		return false
	}

	// download 结束
	if e.opt.eDownload.Activate() {
		return false
	}

	// scheduler 结束
	if e.opt.eScheduler.Activate() {
		return false
	}

	if e.opt.eOutPutter.Activate() {
		return false
	}

	return true
}

func (e *shortEngine) objStream(obj stream.Stream, handle func(to stream.Stream)) {
	switch iter := obj.(type) {
	case *stream.StreamList:
		// 迭代 stream
		iters := iter.Iterator()
		for data := iters(); data != nil; data = iters() {
			handle(data)
		}

	default:
		handle(obj)
	}
}

// 处理来着 spider 的 stream
func (e *shortEngine) fromSpider(to stream.Stream) {
	switch data := to.(type) {
	case *stream.HttpRequest:
		if err := e.opt.eScheduler.Put(data); err != nil {
			wlog.Error(err)
		}
	case stream.Item:
		e.opt.eOutPutter.Put(data)
	}
}

// 处理来自 download 的 stream
func (e *shortEngine) fromDownload(to stream.Stream) {
	switch data := to.(type) {
	case *stream.HttpResponse:

		// 有希望提交出 obj 对象
		item, err := e.opt.eSpider.ParseResp(data)
		if err != nil {
			wlog.Error(err)
			return
		}
		e.objStream(item, e.fromSpider)

	}
}

func (e *shortEngine) fromScheduler(to stream.Stream) {
	switch req := to.(type) {
	case *stream.HttpRequest:

		resp, err := e.opt.eDownload.Do(req)

		if err != nil {
			wlog.Error(err)
			return
		}

		e.fromDownload(resp)
	}
}

func (e *shortEngine) Start(ctx context.Context) {
	// 初始化管道
	if err := e.opt.eOutPutter.OpenPipline(); err != nil {
		wlog.Error(err)
	}

	// 初始化 spider
	e.poll().Submit(func() {
		for _, sp := range e.spiders {
			e.fromSpider(sp.StartRequest())
		}
	})

	// 调度流程

	checkTicker := time.NewTicker(e.opt.checkTime)

	reqCh := e.opt.eScheduler.RecvCtxCancel(ctx)

	for {
		select {
		case <-checkTicker.C:
			if e.canClose() {
				return
			}
		case req := <-reqCh:
			e.poll().Submit(func() {
				e.fromScheduler(req)
			})
		}
	}

}

func (e *shortEngine) Close() {
	if err := e.opt.eScheduler.Close(); err != nil {
		wlog.Error(err)
	}

	if err := e.opt.eOutPutter.Close(); err != nil {
		wlog.Error(err)
	}

	wlog.Info("the spider shuts down successfully")
}

func (e *shortEngine) Register(sp spider.Spider) error {
	if sp == nil {
		return RegisterNilErr
	}
	e.spiders = append(e.spiders, sp)
	return nil
}

func ShortEngine(opts ...optionFunc) *shortEngine {
	ops := loadOption(opts...)

	return &shortEngine{
		opt: ops,
	}
}
