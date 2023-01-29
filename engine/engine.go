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

type DownloadComponent interface {
	SendCop
	ReceiveCop
	io.Closer
}

type SchedulerComponent interface {
	SendCop
	ReceiveCop
	io.Closer
	Finish() error
}

type SpiderComponent interface {
	SendCop
	ReceiveCop
	io.Closer
}
