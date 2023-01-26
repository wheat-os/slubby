package scheduler

import (
	"io"
	"os"
	"sync"

	"github.com/wheat-os/slubby/v2/pkg/cuckoofilter"
)

type CuckooFilter struct {
	mu     sync.Mutex
	cuckoo *cuckoofilter.Filter
	pFile  string

	isClose bool
}

func (c *CuckooFilter) Close() error {
	return nil
}

// PassFilter 检查是否通过过滤器
func (c *CuckooFilter) PassFilter(b []byte) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClose {
		return true
	}

	return !c.cuckoo.AddUnique(b)
}

func loadCuckoo(filePath string) (*cuckoofilter.Filter, error) {
	fs, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fs.Close()

	b, err := io.ReadAll(fs)
	if err != nil {
		return nil, err
	}

	return cuckoofilter.DecodeFrom(b)
}

// WithCuckooFilterP97 使用布谷鸟过滤器，97 准确率
func WithCuckooFilterP97(syPath string) Option {
	cf, err := loadCuckoo(syPath)
	if err != nil {
		cf = cuckoofilter.NewFilter(4, 8, 10000000, cuckoofilter.TableTypePacked)
	}

	return func(s *SlubbyScheduler) {
		s.filter = &CuckooFilter{
			cuckoo: cf,
			pFile:  syPath,
		}
	}
}

// WithCuckooFilterP99 使用布谷鸟过滤器，99.99 准确率,
func WithCuckooFilterP99(syPath string) Option {
	cf, err := loadCuckoo(syPath)
	if err != nil {
		cf = cuckoofilter.NewFilter(4, 16, 10000000, cuckoofilter.TableTypePacked)
	}

	return func(s *SlubbyScheduler) {
		s.filter = &CuckooFilter{
			cuckoo: cf,
			pFile:  syPath,
		}
	}
}

type uselessFilter struct{}

func (n *uselessFilter) Close() error {
	return nil
}

// PassFilter 检查是否通过过滤器
func (n *uselessFilter) PassFilter(b []byte) bool {
	return true
}

// WithUselessFilter 全部放行过滤器
func WithUselessFilter() Option {
	return func(s *SlubbyScheduler) {
		s.filter = &uselessFilter{}
	}
}
