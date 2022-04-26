package bitset

import (
	"bytes"
	"encoding/binary"
)

const (
	bitSize = 8
)

type bit = uint8

type BitSet struct {
	length uint64
	set    []byte
}

func NewBitSet(size uint64) *BitSet {
	return &BitSet{
		length: size,
		set:    make([]byte, (size/bitSize)+1),
	}
}

func NewBitSetBySetContent(set []byte) *BitSet {
	length := binary.BigEndian.Uint64(set[:8])
	return &BitSet{
		length: length,
		set:    set[8:],
	}
}

// 扩容
func (b *BitSet) flashing(size uint64) {
	if size <= b.length || int(size/bitSize)+1 <= len(b.set) {
		return
	}

	dst := make([]byte, (size/bitSize)+1)
	copy(dst, b.set)
	b.length = size
	b.set = dst
}

func (b *BitSet) SetBit(offset uint64, boolen bool) {
	// set true
	i, sv := offset/bitSize, offset%bitSize

	// 尝试扩容
	b.flashing(offset)
	base := bit(1 << (bitSize - 1))

	if boolen {
		b.set[i] |= (base >> sv)
		return
	}

	b.set[i] &= ^(base >> sv)
}

func (b *BitSet) Check(offset uint64) bool {

	if offset > b.length {
		return false
	}

	base := bit(1 << (bitSize - 1))
	i, sv := offset/bitSize, base>>(offset%bitSize)

	return (b.set[i] & sv) != 0
}

func (b *BitSet) Size() uint64 {
	return b.length
}

func (b *BitSet) EncodeBitSet() []byte {
	buf := bytes.NewBuffer(nil)
	bf := make([]byte, 8)
	binary.BigEndian.PutUint64(bf, b.length)
	buf.Write(bf)
	buf.Write(b.set)
	return buf.Bytes()
}

func (b *BitSet) Reset() {
	for i := 0; i < len(b.set); i++ {
		b.set[i] &= 0
	}
}
