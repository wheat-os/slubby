package stream

import (
	"io"
)

type Cover uint8

// Cover 标识来自不同位置的消息
const (
	UnknownCover Cover = iota
	DownloadCover
	SchedulerCover
	SpiderCover
	OutputDeviceCover
)

type FromToCover interface {
	Form() Cover
	To() Cover
	SetForm(c Cover)
	SetTo(c Cover)
}

type TargetCover uint16

// Form 标记 传输 来源
func (t *TargetCover) Form() Cover {
	return Cover(*t >> 8)
}

// To 标记 传输 目标
func (t *TargetCover) To() Cover {
	return Cover(*t << 8 >> 8)
}

// SetForm 设置流传递来源
func (t *TargetCover) SetForm(c Cover) {
	*t = *t<<8>>8 | (TargetCover(c) << 8)
}

// SetTo  设置流传递对象
func (t *TargetCover) SetTo(c Cover) {
	*t = *t>>8<<8 | TargetCover(c)
}

// NewTargetCover 新创建来源目标标记
func NewTargetCover(from, to Cover) TargetCover {
	return TargetCover(from)<<8 | TargetCover(to)
}

type Stream interface {
	io.Closer
	FromToCover
	Context

	// ReplaceCtx 替换或者获取 stream 的上下文对象, 传递为空时不发生替换, 返回旧的上下文
	ReplaceCtx(ctx Context) Context
}

// BackgroundStream 默认 stream
type BackgroundStream struct {
	TargetCover
	Context
}

func (b *BackgroundStream) Close() error {
	return nil
}

func (b *BackgroundStream) MarshalBinary() (data []byte, err error) {
	return nil, nil
}

// ReplaceCtx 上下文替换方案
func (b *BackgroundStream) ReplaceCtx(ctx Context) Context {
	vtr := b.Context
	if ctx != nil {
		b.Context = vtr
	}
	return vtr
}

// Background 获取默认 stream
func Background() Stream {
	return &BackgroundStream{Context: &Meta{}}
}

// Error 携带错误的 stream
func Error(err error) Stream {
	errStream := Background()
	errStream.SetErr(err)
	return errStream
}

// FromError 指定来源错误
func FromError(from Cover, err error) Stream {
	errStream := Error(err)
	errStream.SetForm(from)
	return errStream
}
