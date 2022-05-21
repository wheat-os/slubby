package outputter

import (
	"github.com/panjf2000/ants/v2"
	"github.com/wheat-os/slubby/stream"
)

type Outputter interface {
	Put(item stream.Item)

	Close() error
	Activate() bool

	OpenPipline() error
}

type shortOutputter struct {
	opt *option
}

func (s *shortOutputter) poll() *ants.Pool {
	return s.opt.poll
}

func (s *shortOutputter) Put(item stream.Item) {
	if s.opt.pip == nil {
		return
	}

	s.poll().Submit(func() {
		s.opt.pip.ProcessItem(item)
	})
}

func (s *shortOutputter) Close() error {
	s.poll().Release()

	if s.opt.pip != nil {
		return s.opt.pip.CloseSpider()
	}
	return nil
}

func (s *shortOutputter) OpenPipline() error {
	if s.opt.pip == nil {
		return nil
	}

	return s.opt.pip.OpenSpider()
}

func (s *shortOutputter) Activate() bool {
	return s.poll().Running() != 0
}

func ShortOutputter(opts ...optionFunc) Outputter {
	opt := loadOPtion(opts...)
	return &shortOutputter{
		opt: opt,
	}
}
