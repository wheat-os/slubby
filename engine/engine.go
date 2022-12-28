package engine

import (
	"github.com/wheat-os/slubby/v2/stream"
	"io"
)

// SendCop
//
// Streaming 方法用于基础推流，engine -> cop 的过程，Streaming 是非阻塞方法
type SendCop interface {
	Streaming(data stream.Stream) error
}

// ReceiveCop
//
// ReceiveCop 方法用于基础流接受，用在 engine <- cop 可以向引擎主动推送流
type ReceiveCop interface {
	BackStream() <-chan stream.Stream
}

// SendAndReceiveComponent 推送以及接受组件
type SendAndReceiveComponent interface {
	SendCop
	ReceiveCop
	io.Closer
}

// SendComponent 只推送组件
type SendComponent interface {
	SendCop
	io.Closer
}

// ReceiveComponent 只接受组件
type ReceiveComponent interface {
	ReceiveCop
	io.Closer
}
