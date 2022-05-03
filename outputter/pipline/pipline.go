package pipline

import (
	"gitee.com/wheat-os/slubby/stream"
	"github.com/pkg/errors"
)

type Pipline interface {
	OpenSpider() error

	CloseSpider() error

	ProcessItem(item stream.Item) stream.Item
}

type groupPipline struct {
	piplines []Pipline
}

func (g *groupPipline) OpenSpider() (err error) {

	for _, pip := range g.piplines {
		if pErr := pip.OpenSpider(); pErr != nil {
			err = errors.Wrap(err, pErr.Error())
		}
	}

	return err
}

func (g *groupPipline) CloseSpider() (err error) {
	for _, pip := range g.piplines {
		if pErr := pip.CloseSpider(); pErr != nil {
			err = errors.Wrap(err, pErr.Error())
		}
	}

	return err
}

func (g *groupPipline) ProcessItem(item stream.Item) stream.Item {
	for _, pip := range g.piplines {
		if item == nil {
			break
		}

		item = pip.ProcessItem(item)
	}

	return nil
}

// new pipline group
func GroupPipline(pip ...Pipline) Pipline {
	return &groupPipline{
		piplines: pip,
	}
}
