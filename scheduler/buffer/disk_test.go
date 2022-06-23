package buffer

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/stream"
)

func Test_encodeDiskEntry(t *testing.T) {
	req, err := stream.Request(nil, "www.baidu.com", nil)
	require.NoError(t, err)
	ent, err := newDiskEntryByHttpRequest(req, 0)
	require.NoError(t, err)

	entBy := encodeDiskEntry(ent)
	entTo, err := decodeDiskEntry(entBy)
	require.NoError(t, err)

	require.Equal(t, entTo, ent)

}

func TestFileWrite(t *testing.T) {
	f, err := os.Create("www.test")
	require.NoError(t, err)

	f.WriteString("123")
	f.Seek(0, io.SeekStart)
	f.WriteString("22")
	f.Seek(0, io.SeekStart)
	k, err := ioutil.ReadAll(f)
	require.NoError(t, err)
	require.Equal(t, string(k), "223")
	f.Close()

	os.Remove("www.test")
}

// func TestNewDiskQueue(t *testing.T) {
// 	disk := NewDiskQueue("./.disk.unb")
// 	req, err := stream.Request(nil, "www.baidu.com", nil)
// 	require.NoError(t, err)
// 	err = disk.Put(req)
// 	require.NoError(t, err)
// 	reqI, err := disk.Get()
// 	require.NoError(t, err)
// 	require.Equal(t, req.URL, reqI.URL)

// 	for i := 0; i < 3000; i++ {
// 		req, err := stream.Request(nil, fmt.Sprintf("www.test.com/%d", i), nil)
// 		require.NoError(t, err)
// 		err = disk.Put(req)
// 		require.NoError(t, err)
// 	}

// 	for i := 0; i < 3000; i++ {
// 		req, err := stream.Request(nil, fmt.Sprintf("www.test.com/%d", i), nil)
// 		require.NoError(t, err)
// 		reqI, err = disk.Get()
// 		require.NoError(t, err)

// 		require.Equal(t, req.URL, reqI.URL)
// 	}

// 	_, err = disk.Get()
// 	require.Error(t, err)

// 	disk.Close()

// 	os.Remove("./.disk.unb")

// }

// func Test_diskQueue_encodeHeader(t *testing.T) {
// 	n := diskQueue{
// 		head:      0,
// 		tail:      388,
// 		beginSeek: 256,
// 		length:    2,
// 		factor:    6,
// 	}

// 	p := diskQueue{}

// 	buf := n.encodeHeader()
// 	p.decodeHeader(buf)

// 	require.Equal(t, &n, &p)
// }
