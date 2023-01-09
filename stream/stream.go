package stream

import (
	"encoding"
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
	encoding.BinaryMarshaler
	FromToCover
	Context
}

// BackgroundStream 默认 stream
type BackgroundStream struct {
	TargetCover
	Meta
}

func (b *BackgroundStream) Close() error {
	return nil
}

func (b *BackgroundStream) MarshalBinary() (data []byte, err error) {
	return nil, nil
}

// Background 获取默认 stream
func Background() Stream {
	return &BackgroundStream{}
}

// Error 携带错误的 stream
func Error(err error) Stream {
	return &BackgroundStream{}
}

// FromError 指定来源错误
func FromError(from Cover, err error) Stream {
	errStream := Error(err)
	errStream.SetForm(from)
	return errStream
}
