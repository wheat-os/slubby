package buffer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/stream"
)

func TestQueue_Put(t *testing.T) {
	queue := ShortQueue()
	req, err := stream.Request(nil, "www.baidu.com", nil)
	require.NoError(t, err)
	queue.Put(req)

	require.Equal(t, queue.Size(), 1)

	nReq, _ := queue.Get()

	require.Equal(t, queue.Size(), 0)

	require.Equal(t, nReq, req)

	// 并发测试

	for i := 0; i < 10; i++ {
		go func() {
			queue.Put(req)
		}()
	}

	for i := 0; i < 10; i++ {
		queue.Get()
	}

	require.Equal(t, queue.Size(), 0)

}
