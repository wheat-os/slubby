package pipline

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/stream"
)

type testPipline struct {
	num int
}

func (t *testPipline) OpenSpider() error {
	panic("not implemented") // TODO: Implement
}

func (t *testPipline) CloseSpider() error {
	panic("not implemented") // TODO: Implement
}

func (t *testPipline) ProcessItem(item stream.Item) stream.Item {
	t.num++
	return nil
}

func Test_shuntPipline_ProcessItem(t *testing.T) {
	pip := &testPipline{num: 0}
	shutPip := ShuntPIpline("uid", pip)

	item := stream.BasicItem(nil)
	itemP := stream.BasicItemBandName(nil, "uid")

	shutPip.ProcessItem(item)
	shutPip.ProcessItem(itemP)
	shutPip.ProcessItem(item)

	require.Equal(t, pip.num, 1)
}
