package stream

import (
	"encoding"
	"io"
)

type Cover uint8

// Cover 标识来自不同位置的消息
const (
	DownloadCover Cover = iota
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
}
