package bitset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitSet_SetBit(t *testing.T) {

	bitSet := NewBitSet(300)

	// 不越界测试
	bitSet.SetBit(10, true)
	bitSet.SetBit(300, true)
	bitSet.SetBit(11, true)
	bitSet.SetBit(0, true)

	flag := bitSet.Check(10)
	require.True(t, flag)
	flag = bitSet.Check(300)
	require.True(t, flag)
	flag = bitSet.Check(11)
	require.True(t, flag)
	flag = bitSet.Check(0)
	require.True(t, flag)
	// 修改 其他 bit 不应该影响其他
	bitSet.SetBit(10, false)

	flag = bitSet.Check(10)
	require.False(t, flag)
	flag = bitSet.Check(300)
	require.True(t, flag)
	flag = bitSet.Check(11)
	require.True(t, flag)
	flag = bitSet.Check(0)
	require.True(t, flag)

	require.Equal(t, bitSet.Size(), uint(300))

	// 越界测试
	bitSet.SetBit(1000, true)
	flag = bitSet.Check(1000)
	require.True(t, flag)

	flag = bitSet.Check(20000)
	require.False(t, flag)

	require.Equal(t, bitSet.Size(), uint(1000))
}

func TestNewBitSetBySetContent(t *testing.T) {
	bitset := NewBitSet(100)
	bitset.SetBit(10, true)
	bitset = NewBitSetBySetContent(bitset.Bytes())

	f := bitset.Check(10)
	require.True(t, f)
}
