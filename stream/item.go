package stream

import "encoding/binary"

type Item struct {
	TargetCover
}

// Close 关闭 item， 默认为空
func (i *Item) Close() error {
	return nil
}

// MarshalBinary 编码 item， 默认编码 cover
func (i *Item) MarshalBinary() (data []byte, err error) {
	buf := make([]byte, 0, 2)
	buf = binary.BigEndian.AppendUint16(buf, uint16(i.TargetCover))
	return buf, err
}
