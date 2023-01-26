package buffer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListBuffer_Get(t *testing.T) {
	buf := NewListBuffer(10)
	val, err := buf.Get()
	require.Nil(t, val)
	require.Equal(t, err, ErrBufferIsEmpty)

	tv := "123"
	err = buf.Put(tv)
	require.NoError(t, err)

	val, err = buf.Get()
	require.NoError(t, err)
	require.Equal(t, val, tv)

	val, err = buf.Get()
	require.Nil(t, val)
	require.Equal(t, err, ErrBufferIsEmpty)
}
