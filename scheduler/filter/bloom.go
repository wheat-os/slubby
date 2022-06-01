package filter

import (
	"io/ioutil"
	"math"
	"os"
	"sync"

	"github.com/wheat-os/slubby/pkg/bitset"
	"github.com/wheat-os/slubby/stream"
)

type Filter interface {
	// 是否通过 过滤器
	Passage(req *stream.HttpRequest) (bool, error)

	// 重置过滤器
	Reset()

	Close() error
}

type BloomFilter struct {
	set       bitset.BitSet
	mu        sync.Mutex
	filePath  string
	hashCount int
}

func (b *BloomFilter) hash(content []byte, seed uint64) uint64 {
	var result uint64

	for i := 0; i < len(content); i++ {
		result = result + seed*uint64(content[i])
	}

	return result % b.set.Size()
}

// 是否通过 过滤器
func (b *BloomFilter) Passage(req *stream.HttpRequest) (bool, error) {

	b.mu.Lock()
	defer b.mu.Unlock()

	buf, err := stream.EncodeHttpRequest(req)
	if err != nil {
		return false, err
	}

	var isPassage bool
	for seed := 1; seed <= b.hashCount; seed++ {
		hashVal := b.hash(buf, uint64(seed*1024))

		// 有一个 val 不存在通过过滤器
		if !b.set.Check(hashVal) {
			isPassage = true
			b.set.SetBit(hashVal, true)
		}
	}

	return isPassage, nil
}

// 重置过滤器
func (b *BloomFilter) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.set.Reset()
}

func (b *BloomFilter) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.filePath == "" {
		return nil
	}

	f, err := os.Create(b.filePath)
	if err != nil {
		return err
	}

	defer f.Close()
	f.Write(b.set.EncodeBitSet())
	return nil
}

func loadBloomSet(size uint64, filePath string) bitset.BitSet {
	if filePath != "" {
		if f, err := os.Open(filePath); err == nil {
			defer f.Close()
			if buf, err := ioutil.ReadAll(f); err == nil {
				return bitset.NewBitSetBySetContent(buf)
			}
		}
	}

	return bitset.NewBitSet(size)
}

// 根据数量级和损失率生成 bloom filter
// magnitude 插入数据数量级
// loss 损失率
func NewBloomFilter(magnitude int, loss float64, filePath string) Filter {
	size := -((math.Log(loss) * float64(magnitude)) / math.Log(2) / math.Log(2))
	hashCount := (size / float64(magnitude)) * math.Log(2)

	size = math.Ceil(size)
	hashCount = math.Floor(hashCount + 0.5)

	return &BloomFilter{
		set:       loadBloomSet(uint64(size), filePath),
		hashCount: int(hashCount),
		filePath:  filePath,
	}
}

func ShortBloomFilter() Filter {
	const (
		magnitude = 100000
		loss      = 0.05
	)
	return NewBloomFilter(magnitude, loss, "")
}
