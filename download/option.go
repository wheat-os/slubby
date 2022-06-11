package download

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wheat-os/slubby/download/limiter"
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/pkg/sundry"
)

type option struct {
	// 下载器最大并发数
	concurrentRequest int
	// 每个 FQDN 下载完成后的延迟时间
	limiter limiter.Limiter

	client *http.Client
	poll   *ants.Pool
	once   sync.Once

	// mid
	middle.Middleware

	// retry
	retry int
}

// 初始化 poll
func (o *option) Poll() *ants.Pool {
	o.once.Do(func() {
		poll, err := ants.NewPool(
			o.concurrentRequest,
			ants.WithPanicHandler(sundry.AntsWlogHandlePanic),
		)
		if err != nil {
			panic(err)
		}
		o.poll = poll
	})

	return o.poll
}

func loadOption(opts ...optionFunc) *option {
	// init option
	ops := &option{
		//  -1 == Max
		concurrentRequest: 8,
		client:            http.DefaultClient,
		retry:             2,
	}

	for _, opt := range opts {
		opt(ops)
	}

	ops.Poll()
	return ops
}

type optionFunc = func(opt *option)

func WithConcurrentRequest(num int) optionFunc {
	return func(opt *option) {
		opt.concurrentRequest = num
	}
}

func WithLimiter(lim limiter.Limiter) optionFunc {
	return func(opt *option) {
		opt.limiter = lim
	}
}

func WithClient(c *http.Client) optionFunc {
	return func(opt *option) {
		opt.client = c
	}
}

func WithTimeout(timeout time.Duration) optionFunc {
	return func(opt *option) {
		transport := http.DefaultTransport.(*http.Transport)

		transport.DialContext = (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}
}

func WithDownloadMiddle(mid middle.Middleware) optionFunc {
	return func(opt *option) {
		opt.Middleware = mid
	}
}

func WithDownloadRetry(retry int) optionFunc {
	if retry <= 0 {
		retry = 2
	}
	return func(opt *option) {
		opt.retry = retry
	}
}
