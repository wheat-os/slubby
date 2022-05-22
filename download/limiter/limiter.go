package limiter

import (
	"sync"
	"time"

	"github.com/wheat-os/slubby/stream"
)

type Limiter interface {
	Allow(data stream.Stream) bool
}

type fqdnSetting func() (string, time.Duration)

func SetFqdnLimiter(fqdn string, delay time.Duration) fqdnSetting {
	return func() (string, time.Duration) {
		return fqdn, delay
	}
}

type shortLimiter struct {
	bucket map[string]chan struct{}
	mu     sync.Mutex

	fqdnSet map[string]time.Duration
	delay   time.Duration
}

func (s *shortLimiter) runLimiter(fqdn string) {

	ch := make(chan struct{})
	s.bucket[fqdn] = ch

	go func() {
		delay := s.delay
		if fqdnDelay, ok := s.fqdnSet[fqdn]; ok {
			delay = fqdnDelay
		}

		ticket := time.NewTicker(delay)
		for {
			<-ticket.C
			ch <- struct{}{}
		}
	}()
}

func (s *shortLimiter) Allow(data stream.Stream) bool {

	if data == nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ch, ok := s.bucket[data.FQDN()]
	// 不存在 fqdn 限流器
	if !ok {
		s.runLimiter(data.FQDN())
		return true
	}

	<-ch
	return true
}

func ShortLimiter(defaultDelay time.Duration, set ...fqdnSetting) Limiter {

	fqdnSet := make(map[string]time.Duration)

	for _, s := range set {
		fqdn, delay := s()
		fqdnSet[fqdn] = delay
	}

	return &shortLimiter{
		bucket:  make(map[string]chan struct{}),
		delay:   defaultDelay,
		fqdnSet: fqdnSet,
	}
}
