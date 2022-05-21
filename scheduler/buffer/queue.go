package buffer

import (
	"math"
	"sync"

	"github.com/wheat-os/slubby/stream"
)

type Buffer interface {
	// 请求逻辑
	Put(req *stream.HttpRequest) error
	Get() (*stream.HttpRequest, error)

	Size() int
	Cap() int

	Close() error
}

// CAS 队列
type Queue struct {
	head   *qNode
	tail   *qNode
	length int
	cond   *sync.Cond
}

type qNode struct {
	Data *stream.HttpRequest
	Next *qNode
}

func ShortQueue() Buffer {
	basic := &qNode{}

	return &Queue{
		head:   basic,
		tail:   basic,
		length: 0,
		cond:   sync.NewCond(&sync.Mutex{}),
	}
}

// 请求逻辑
func (q *Queue) Put(req *stream.HttpRequest) error {
	node := &qNode{Data: req}
	q.cond.L.Lock()

	q.tail.Next = node
	q.tail = node
	q.length += 1

	q.cond.L.Unlock()
	q.cond.Signal()

	return nil
}

func (q *Queue) Get() (*stream.HttpRequest, error) {
	q.cond.L.Lock()
	for q.length == 0 {
		q.cond.Wait()
	}

	node := q.head.Next
	q.head.Next = node.Next

	q.length -= 1

	if q.length == 0 {
		q.tail = q.head
	}

	q.cond.L.Unlock()

	return node.Data, nil

}

func (q *Queue) Size() int {
	return q.length
}

// 无穷大
func (q *Queue) Cap() int {
	return math.MaxInt
}

func (q *Queue) Close() error {
	return nil
}
