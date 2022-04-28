package scheduler

import (
	"gitee.com/wheat-os/slubby/stream"
	"github.com/pkg/errors"
)

type Scheduler interface {
	Put(req *stream.HttpRequest) error
	Get() (*stream.HttpRequest, error)

	Close() error
	Activate() bool
}

type shortScheduler struct {
	opt *option
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

func (s *shortScheduler) Activate() bool {
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
func ShortScheduler() Scheduler {
	opt := loadOption()

	return &shortScheduler{
		opt: opt,
	}
}

func NewScheduler(opts ...optionFunc) Scheduler {
	opt := loadOption(opts...)

	return &shortScheduler{
		opt: opt,
	}
}
