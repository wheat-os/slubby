package download

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type option struct {
	// 下载器最大并发数
	concurrentRequest int
	// 每个 FQDN 下载完成后的延迟时间
	delay int

	client *http.Client
	poll   *ants.Pool
	once   sync.Once
}

// 初始化 poll
func (o *option) Poll() *ants.Pool {
	o.once.Do(func() {
		poll, err := ants.NewPool(o.concurrentRequest)
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
		delay:             -1,
		client:            http.DefaultClient,
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

func WithFQDNDelay(delay int) optionFunc {
	return func(opt *option) {
		opt.delay = delay
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
