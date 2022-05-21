package download

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/stream"
)

func TestDownload(t *testing.T) {
	pool, err := ants.NewPool(1000)
	require.NoError(t, err)

	require.Equal(t, pool.Free(), 1000)
	require.Equal(t, pool.Running(), 0)

	pool.Submit(func() {
		fmt.Println(1)
	})

}

func Test_shortDownload_Do(t *testing.T) {
	download := ShortDownload()

	req, err := stream.Request(nil, "https://www.baidu.com", nil)
	require.NoError(t, err)

	resp, err := download.Do(req)
	require.NoError(t, err)

	require.Equal(t, resp.StatusCode, 200)

	// 并发测试
	wait := sync.WaitGroup{}

	wait.Add(20)

	for i := 0; i < 20; i++ {
		go func() {
			_, err := download.Do(req)
			require.NoError(t, err)
			wait.Done()
		}()
	}

	wait.Wait()
}

type testMiddle struct {
}

func (t *testMiddle) BeforeDownload(m *middle.M, req *stream.HttpRequest) (*stream.HttpRequest, error) {
	fmt.Println("before")
	return req, nil
}

func (t *testMiddle) AfterDownload(m *middle.M, req *stream.HttpRequest, resp *stream.HttpResponse) (*stream.HttpResponse, error) {
	fmt.Println("after")
	return resp, nil
}

func (t *testMiddle) ProcessErr(m *middle.M, req *stream.HttpRequest, err error) {
	panic("err")
}

func Test_shortDownload_Do_Middle(t *testing.T) {
	download := ShortDownload(
		WithFQDNDelay(20),
		WithTimeout(20*time.Second),
		WithDownloadMiddle(&testMiddle{}),
	)

	req, err := stream.Request(nil, "https://www.baidu.com", nil)
	require.NoError(t, err)

	resp, err := download.Do(req)
	require.NoError(t, err)

	fmt.Println(resp)
}
