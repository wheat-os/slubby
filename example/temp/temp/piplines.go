package temp

import (
	"gitee.com/wheat-os/slubby/stream"
)

type TempPipline struct{}

func (t *TempPipline) OpenSpider() error {
	return nil
}

func (t *TempPipline) CloseSpider() error {
	return nil
}

func (t *TempPipline) ProcessItem(item stream.Item) stream.Item {
	return item
}
