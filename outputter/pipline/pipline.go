package pipline

import (
	"github.com/pkg/errors"
	"github.com/wheat-os/slubby/stream"
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

type shuntPipline struct {
	pip      Pipline
	itemName string
}

func (s *shuntPipline) OpenSpider() error {
	return s.pip.OpenSpider()
}

func (s *shuntPipline) CloseSpider() error {
	return s.pip.CloseSpider()
}

func (s *shuntPipline) ProcessItem(item stream.Item) stream.Item {
	if s.itemName != item.IName() {
		return item
	}

	return s.pip.ProcessItem(item)
}

func ShuntPIpline(iName string, pip Pipline) Pipline {
	return &shuntPipline{
		pip:      pip,
		itemName: iName,
	}
}
