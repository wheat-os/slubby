package stream

import (
	"testing"

	perr "gitee.com/wheat-os/slubby/pkg/error"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestEncodeShortStream(t *testing.T) {
	shortInfo := SpiderInfo("spider1", "www.baidu.com")

	buf := EncodeShortStream(shortInfo)

	spiderInfo, err := DecodeShortStream(buf)
	require.NoError(t, err)

	require.Equal(t, shortInfo, spiderInfo)

	buf = []byte("awdwadhibbawd")
	_, err = DecodeShortStream(buf)
	require.Error(t, err)

	require.True(t, errors.Is(err, perr.InvalidEncodingErr))
}

func TestMustRegisterSpiderStram(t *testing.T) {
	spider := &TestSpider{
		Stream: SpiderInfo("test", "www.qq.com"),
	}

	callback := GetCallbackFuncByName(spider, "Parse")
	require.Nil(t, callback)

	MustRegisterSpiderStram(spider, spider.Parse, spider.GetList)

	callback = GetCallbackFuncByName(spider, "Parse")
	require.NotNil(t, callback)
	callback(nil)

	callback = GetCallbackFuncByName(spider, "ToList")
	require.Nil(t, callback)

	MustRegisterSpiderStram(spider, spider.ToList)

	callback = GetCallbackFuncByName(spider, "ToList")
	require.NotNil(t, callback)
}
