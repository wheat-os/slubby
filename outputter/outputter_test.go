package outputter

import (
	"testing"

	"gitee.com/wheat-os/slubby/stream"
	"gitee.com/wheat-os/wlog"
	"github.com/stretchr/testify/require"
)

type testPipline struct{}

func (t *testPipline) OpenSpider() error {
	wlog.Info("open in pipline")
	return nil
}

func (t *testPipline) CloseSpider() error {
	wlog.Info("close pipline")
	return nil
}

func (t *testPipline) ProcessItem(item stream.Item) stream.Item {
	wlog.Info(item)
	return nil
}

func TestShortOutputter(t *testing.T) {
	ts := &testPipline{}

	opt := ShortOutputter(
		WithPipline(ts),
	)

	require.Equal(t, opt.Activate(), false)

	err := opt.OpenPipline()
	require.NoError(t, err)

	opt.Put(nil)

	for opt.Activate() {
	}

	err = opt.Close()
	require.NoError(t, err)
}
