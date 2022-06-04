package stream

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamLists(t *testing.T) {
	ts := &TestSpider{}
	url := "www.baidu.com"
	url2 := "www.qq.com"

	streamValues := make([]Stream, 0)
	req, err := Request(ts, url, ts.GetList)
	require.NoError(t, err)
	streamValues = append(streamValues, req)

	req, err = Request(ts, url2, ts.GetList)
	require.NoError(t, err)
	streamValues = append(streamValues, req)

	iter := StreamLists(ts, streamValues...).(*StreamList).Iterator()

	nReq := iter()
	require.Equal(t, nReq, streamValues[0])
	nReq = iter()
	require.Equal(t, nReq, streamValues[1])
	nReq = iter()
	require.Nil(t, nReq)

}

func TestStreamListRangeInt(t *testing.T) {
	sts, err := StreamListRangeInt(nil, func(i int) (Stream, error) {
		return Request(nil, "http://www.test/com", nil)
	}, 10, 20)

	require.NoError(t, err)

	require.Equal(t, len(sts.(*StreamList).streams), 11)
}

func TestStreamListRangeString(t *testing.T) {
	sts, err := StreamListRangeString(nil, func(i string) (Stream, error) {
		return Request(nil, "http://www.test/com", nil)
	}, []string{"1", "2", "3"})

	require.NoError(t, err)

	require.Equal(t, len(sts.(*StreamList).streams), 3)
}

func TestStreamListRangeFloat(t *testing.T) {
	testFloat := []float64{1.1, 1.2, 1.3, 1.333}
	sts, err := StreamListRangeFloat(nil, func(value float64) (Stream, error) {
		req, err := Request(nil, "http://www.baidu.com", nil)
		if err != nil {
			return nil, err
		}

		req.Meta["val"] = value
		return req, nil
	}, testFloat)

	require.NoError(t, err)

	iter := sts.(*StreamList).Iterator()

	for _, v := range testFloat {
		it := iter()
		value := it.(*HttpRequest).Meta["val"]
		require.Equal(t, v, value)
	}
}
