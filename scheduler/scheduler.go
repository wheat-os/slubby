package scheduler

import (
	"context"
	"sync/atomic"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/wheat-os/slubby/stream"
	"github.com/wheat-os/wlog"
)

type Scheduler interface {
	Put(req *stream.HttpRequest) error
	Get() (*stream.HttpRequest, error)

	// 引擎通过 cancel 来检查及时退出
	RecvCtxCancel(ctx context.Context) <-chan *stream.HttpRequest

	Close() error
	Activate() bool
}

type shortScheduler struct {
	opt *option

	isCancel bool
}

func (s *shortScheduler) Put(req *stream.HttpRequest) error {
	if s.opt.filter != nil {
		if b, err := s.opt.filter.Passage(req); err != nil || !b {
			return err
		}
	}

	return s.opt.buffer.Put(req)
}

func (s *shortScheduler) Get() (*stream.HttpRequest, error) { return s.opt.buffer.Get() }

// RecvCtxCancel
func (s *shortScheduler) RecvCtxCancel(ctx context.Context) <-chan *stream.HttpRequest {

	ch := make(chan *stream.HttpRequest)
	go func() {
	F:
		for {
			select {
			case <-ctx.Done():
				break F
			default:
				if s.opt.buffer.Size() == 0 {
					continue F
				}

				if req, err := s.Get(); err == nil {
					ch <- req
				} else {
					wlog.Error(err)
				}
			}
		}

		// 关闭标识
		res := true
		atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&s.isCancel)), unsafe.Pointer(&res))
	}()

	return ch
}

func (s *shortScheduler) Activate() bool {
	if s.isCancel {
		return false
	}

	return s.opt.buffer.Size() != 0
}

func (s *shortScheduler) Close() error {
	err := s.opt.buffer.Close()

	if s.opt.filter != nil {
		if filterErr := s.opt.filter.Close(); filterErr != nil {
			err = errors.Wrap(err, filterErr.Error())
		}
	}

	return err
}

// default scheduler
// bloom filter(10000, 0.95)
// queue
func ShortScheduler(opts ...optionFunc) Scheduler {
	opt := loadOption(opts...)

	return &shortScheduler{
		opt:      opt,
		isCancel: false,
	}
}
