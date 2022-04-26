package bitset

const (
	bitSize = 8
)

type bit = uint8

type BitSet struct {
	length uint
	set    []byte
}

func NewBitSet(size uint) *BitSet {
	return &BitSet{
		length: size,
		set:    make([]byte, (size/bitSize)+1),
	}
}

func NewBitSetBySetContent(set []byte) *BitSet {
	return &BitSet{
		length: uint(len(set) * bitSize),
		set:    set,
	}
}

// 扩容
func (b *BitSet) flashing(size uint) {
	if size <= b.length || int(size/bitSize)+1 <= len(b.set) {
		return
	}

	dst := make([]byte, (size/bitSize)+1)
	copy(dst, b.set)
	b.length = size
	b.set = dst
}

func (b *BitSet) SetBit(offset uint, boolen bool) {
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

func (b *BitSet) Check(offset uint) bool {

	if offset > b.length {
		return false
	}

	base := bit(1 << (bitSize - 1))
	i, sv := offset/bitSize, base>>(offset%bitSize)

	return (b.set[i] & sv) != 0
}

func (b *BitSet) Size() uint {
	return b.length
}

func (b *BitSet) Bytes() []byte {
	return b.set
}

func (b *BitSet) Reset() {
	for i := 0; i < len(b.set); i++ {
		b.set[i] &= 0
	}
}
