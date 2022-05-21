package filter

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/stream"
)

func TestShort_BloomFilter(t *testing.T) {

	// 默认 short 为:
	// 10 w 负载
	// 0.05 损失率
	// 计算 size = 623523
	// hash func count = 4
	bloom := ShortBloomFilter().(*BloomFilter)
	require.Equal(t, bloom.hashCount, 4)
	require.Equal(t, bloom.set.Size(), uint(623523))
}

func TestBloomFilter_Passage(t *testing.T) {
	bloom := ShortBloomFilter()

	req, err := stream.Request(nil, "http://www.baidu.com", nil)
	require.NoError(t, err)

	flag, err := bloom.Passage(req)
	require.NoError(t, err)
	require.True(t, flag)

	// new display req
	nReq, err := stream.Request(nil, "http://www.baidu.com", nil)
	require.NoError(t, err)

	flag, err = bloom.Passage(nReq)
	require.NoError(t, err)
	require.False(t, flag)

	nReq, err = stream.Request(nil, "http://www.qq.com", nil)
	require.NoError(t, err)

	flag, err = bloom.Passage(nReq)
	require.NoError(t, err)
	require.True(t, flag)

	// header and body req
	body := `{"data": 123}`
	bodyIo := bytes.NewBufferString(body)
	hReq, err := stream.BodyRequest(nil, http.MethodPost, "http://www.test.com", bodyIo, nil)
	require.NoError(t, err)

	hReq.Header.Add("content-type", "applocation/json")

	flag, err = bloom.Passage(hReq)
	require.NoError(t, err)
	require.True(t, flag)

	// header and body req
	body = `{"data": 123}`
	bodyIo = bytes.NewBufferString(body)
	hReq, err = stream.BodyRequest(nil, http.MethodPost, "http://www.test.com", bodyIo, nil)
	require.NoError(t, err)
	hReq.Header.Add("content-type", "applocation/json")

	// 触发过滤
	flag, err = bloom.Passage(hReq)
	require.NoError(t, err)
	require.False(t, flag)

}

func TestBloomFilter_Close(t *testing.T) {
	os.Remove("./bloom.odb")

	bloom := NewBloomFilter(10000, 0.03, "./bloom.odb")
	req, err := stream.Request(nil, "http://www.baidu.com", nil)
	require.NoError(t, err)
	flag, err := bloom.Passage(req)
	require.NoError(t, err)
	require.True(t, flag)

	err = bloom.Close()
	require.NoError(t, err)

	bloom = NewBloomFilter(10000, 0.03, "./bloom.odb")
	flag, err = bloom.Passage(req)
	require.NoError(t, err)
	require.False(t, flag)

	err = bloom.Close()
	require.NoError(t, err)

}
