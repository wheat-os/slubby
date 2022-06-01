package bitset

import (
	"bytes"
	"encoding/binary"
)

const (
	bitSize = 8
)

type bit = uint8

type bitSet struct {
	length uint64
	set    []byte
}

type BitSet interface {
	SetBit(offset uint64, boolen bool)
	Check(offset uint64) bool
	Size() uint64

	EncodeBitSet() []byte

	Reset()
}

func NewBitSet(size uint64) BitSet {
	return &bitSet{
		length: size,
		set:    make([]byte, (size/bitSize)+1),
	}
}

func NewBitSetBySetContent(set []byte) BitSet {
	length := binary.BigEndian.Uint64(set[:8])
	return &bitSet{
		length: length,
		set:    set[8:],
	}
}

// 扩容
func (b *bitSet) flashing(size uint64) {
	if size <= b.length || int(size/bitSize)+1 <= len(b.set) {
		return
	}

	dst := make([]byte, (size/bitSize)+1)
	copy(dst, b.set)
	b.length = size
	b.set = dst
}

func (b *bitSet) SetBit(offset uint64, boolen bool) {
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

func (b *bitSet) Check(offset uint64) bool {

	if offset > b.length {
		return false
	}

	base := bit(1 << (bitSize - 1))
	i, sv := offset/bitSize, base>>(offset%bitSize)

	return (b.set[i] & sv) != 0
}

func (b *bitSet) Size() uint64 {
	return b.length
}

func (b *bitSet) EncodeBitSet() []byte {
	buf := bytes.NewBuffer(nil)
	bf := make([]byte, 8)
	binary.BigEndian.PutUint64(bf, b.length)
	buf.Write(bf)
	buf.Write(b.set)
	return buf.Bytes()
}

func (b *bitSet) Reset() {
	for i := 0; i < len(b.set); i++ {
		b.set[i] &= 0
	}
}
