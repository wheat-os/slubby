package stream

import "errors"

var (
	ErrUnknownEncoderStreamType = errors.New("unrecognized stream, parse stream type")
)

// Encoder 流编解码器
type Encoder interface {
	// StreamEncode 实现流编码方法
	StreamEncode(stm Stream) ([]byte, error)
	// StreamDecode 重分离的 binary 中解析流信息
	StreamDecode(b []byte) (Stream, error)
}

type BackGroundEncoder struct {
}

func (b BackGroundEncoder) StreamEncode(stm Stream) ([]byte, error) {
	return nil, nil
}

func (b BackGroundEncoder) StreamDecode(p []byte) (Stream, error) {
	return Background(), nil
}

func NewBackGroundEncoder() BackGroundEncoder {
	return BackGroundEncoder{}
}
