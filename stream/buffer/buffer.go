package buffer

import (
	"errors"
	"github.com/wheat-os/slubby/v2/pkg/buffer"
	"github.com/wheat-os/slubby/v2/stream"
)

var (
	ErrStreamBufferIsEmpty    = errors.New("this stream buffer is empty")
	ErrStreamBufferIsFull     = errors.New("this stream buffer is full")
	ErrStreamAssertionFailure = errors.New("this entry parse to stream err")
)

// StreamBuffer 下载器缓冲区组件，应该保证线程安全
type StreamBuffer interface {
	Len() int
	Cap() int
	PutStream(data stream.Stream) error
	GetStream() (stream.Stream, error)
}

type ListBuffer struct {
	list *buffer.ListBuffer
}

func (l *ListBuffer) Len() int {
	return l.list.Len()
}

func (l *ListBuffer) Cap() int {
	return l.list.Cap()
}

func (l *ListBuffer) PutStream(data stream.Stream) error {
	err := l.list.Put(data)

	switch err {
	case nil:
		return nil
	case buffer.ErrBufferIsEmpty:
		return ErrStreamBufferIsEmpty
	default:
		return err
	}
}

func (l *ListBuffer) GetStream() (stream.Stream, error) {
	entry, err := l.list.Get()
	stm, ok := entry.(stream.Stream)
	if !ok {
		return nil, ErrStreamAssertionFailure
	}

	switch err {
	case nil:
		return stm, nil
	case buffer.ErrBufferIsFull:
		return nil, ErrStreamBufferIsFull
	default:
		return stm, err
	}
}

// NewListBuffer 创建一个无等待即使返回缓冲区
func NewListBuffer(cap int) *ListBuffer {
	return &ListBuffer{list: buffer.NewListBuffer(cap)}
}

type ChanBuffer chan stream.Stream

func (c ChanBuffer) Len() int {
	return len(c)
}

func (c ChanBuffer) Cap() int {
	return cap(c)
}

func (c ChanBuffer) PutStream(data stream.Stream) error {
	c <- data
	return nil
}

func (c ChanBuffer) GetStream() (stream.Stream, error) {
	stm := <-c
	return stm, nil
}

// NewChanBuffer 创建使用 go chan 实现的 stream 缓冲区
func NewChanBuffer(cap int) ChanBuffer {
	return make(ChanBuffer, cap)
}
