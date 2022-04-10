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
