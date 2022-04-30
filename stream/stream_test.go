package stream

import (
	"testing"

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

	require.True(t, errors.Is(err, InvalidEncodingErr))
}
