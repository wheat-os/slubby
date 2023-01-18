package buffer

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	ErrBufferIsEmpty = errors.New("this buffer is empty")
	ErrBufferIsFull  = errors.New("this buffer is full")
	ErrBufferTimeout = errors.New("buffer option is timeout")
)

type ListBuffer struct {
	buf  *list.List
	cond *sync.Cond
	cap  int
}

func NewListBuffer(cap int) *ListBuffer {
	return &ListBuffer{buf: list.New(), cond: sync.NewCond(&sync.Mutex{}), cap: cap}
}

// Len 数据长度
func (l *ListBuffer) Len() int {
	return l.buf.Len()
}

// Cap 容量
func (l *ListBuffer) Cap() int {
	return l.cap
}

// Get 无阻塞获取数据
func (l *ListBuffer) Get() (any, error) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	if l.buf.Len() <= 0 {
		return nil, ErrBufferIsEmpty
	}

	ele := l.buf.Front()
	l.buf.Remove(ele)
	return ele.Value, nil
}

// Put 无阻塞推送
func (l *ListBuffer) Put(val any) error {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	if l.Len() >= l.Cap() {
		return ErrBufferIsFull
	}

	l.buf.PushBack(val)
	return nil
}

// WaitPut 阻塞推送
func (l *ListBuffer) WaitPut(val any) error {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()
	// 等待
	for l.Len() >= l.Cap() {
		l.cond.Signal()
		l.cond.Wait()
	}

	l.buf.PushBack(val)
	return nil
}

// WaitGet 阻塞读取
func (l *ListBuffer) WaitGet() (any, error) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()
	for l.Len() == 0 {
		l.cond.Signal()
		l.cond.Wait()
	}

	ele := l.buf.Front()
	l.buf.Remove(ele)
	return ele.Value, nil
}

// TimeoutGet 超时读取方法
func (l *ListBuffer) TimeoutGet(timeout time.Duration) (any, error) {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()
	tk := time.NewTimer(timeout)
	defer tk.Stop()

	for l.Len() == 0 {
		select {
		case <-tk.C:
			return nil, ErrBufferTimeout
		default:
		}

		l.cond.Signal()
		l.cond.Wait()
	}

	ele := l.buf.Front()
	l.buf.Remove(ele)
	return ele.Value, nil
}

// TimeoutPut 超时推送
func (l *ListBuffer) TimeoutPut(val any, timeout time.Duration) error {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()

	// 等待
	tk := time.NewTimer(timeout)
	defer tk.Stop()

	for l.Len() >= l.Cap() {
		select {
		case <-tk.C:
			return ErrBufferTimeout
		default:
		}

		l.cond.Signal()
		l.cond.Wait()
	}

	l.buf.PushBack(val)

	// 唤醒全部消费者
	if l.Len() >= l.Cap()/2 {
		l.cond.Broadcast()
	}
	return nil
}
