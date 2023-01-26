package download

import (
	"github.com/wheat-os/slubby/v2/stream/buffer"
	"github.com/wheat-os/wlog"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/v2/engine"
	"github.com/wheat-os/slubby/v2/stream"
	"github.com/wheat-os/slubby/v2/stream/http"
)

// 实现基础的 slubby download 测试
func TestSlubbyDownload(t *testing.T) {
	tests := []struct {
		name      string
		component engine.SendAndReceiveComponent
	}{
		{
			name:      "基础 slubby download",
			component: NewSlubbyDownload(),
		},
		{
			name: "使用 无阻塞缓冲队列 download",
			component: NewSlubbyDownload(
				WithDownloadBuffer(buffer.NewListBuffer(10)),
				WithTokenBucketLimit(2*time.Second, 1),
			),
		},
		{
			name: "单线程模型",
			component: NewSlubbyDownload(
				WithDownloadBuffer(buffer.NewListBuffer(10)),
				WithTokenBucketLimit(2*time.Second, 1),
				WithDownloadProcess(1),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.component.BackStream()

			// 百度基本测试
			req, err := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
			req.SetForm(stream.SchedulerCover)
			req.SetTo(stream.DownloadCover)
			require.NoError(t, err)
			tt.component.Streaming(req)
			data := <-ch
			require.NoError(t, data.Err())
			resp := data.(*http.StreamResponse)
			b, err := io.ReadAll(resp)
			require.Greater(t, len(b), 0)
			require.NoError(t, err)
			require.Equal(t, resp.To(), stream.SpiderCover)

			// douban 测试
			req, err = http.NewRequest(http.MethodGet, "https://movie.douban.com/top250?start=0&filter=", nil)
			req.Request.Header.Set("user-agent", "Mozilla/5.0")
			req.SetForm(stream.SchedulerCover)
			req.SetTo(stream.DownloadCover)
			require.NoError(t, err)
			tt.component.Streaming(req)
			data = <-ch
			require.NoError(t, data.Err())
			resp = data.(*http.StreamResponse)
			b, err = io.ReadAll(resp)
			require.Greater(t, len(b), 0)
			require.NoError(t, err)
			require.Equal(t, resp.To(), stream.SpiderCover)

			go func() {
				for {
					stm := <-ch
					wlog.Errorf("%t", stm)
				}
			}()

			tt.component.Close()
		})
	}
}
