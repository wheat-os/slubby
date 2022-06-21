package buffer

import (
	"bytes"
	"encoding/binary"
	"hash/crc64"
	"io"
	"math"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/wheat-os/slubby/pkg/tools"
	"github.com/wheat-os/slubby/stream"
	"github.com/wheat-os/wlog"
)

var (
	parseEntryErr = errors.New("resolve disk queue error, disk queue file has been corrupted")
	buildNilEntry = errors.New("cannot code an empty request to entry")

	verifiedReqErr    = errors.New("verified err, disk file has been corrupted or file seek err")
	diskQueueCloseErr = errors.New("the disk queue has been closed and cannot be read or written")
)

const (
	verifiedSize = 8
)

type entry struct {
	length  uint64
	content []byte
	ver     []byte
}

func (e *entry) size() int {
	return 8 + len(e.content) + len(e.ver)
}

func newEntry(req *stream.HttpRequest) (*entry, error) {
	if req == nil {
		return nil, buildNilEntry
	}

	content, err := stream.EncodeHttpRequest(req)
	if err != nil {
		return nil, err
	}

	return &entry{
		length:  uint64(len(content)),
		content: content,
	}, nil
}

func encodeDiskEntry(e *entry) []byte {
	buf := bytes.NewBuffer(nil)

	// build length
	uIntBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(uIntBuf, e.length)
	buf.Write(uIntBuf)

	// build content
	buf.Write(e.content)

	// build verified hash
	crc := crc64.New(crc64.MakeTable(verifiedSize))
	crc.Write(e.content)
	ver := crc.Sum(nil)
	buf.Write(ver)

	return buf.Bytes()
}

func getDIskEntry(buf io.Reader) (*entry, error) {
	uIntBuf := make([]byte, 8)

	// decode length
	n, err := buf.Read(uIntBuf)
	if err != nil {
		return nil, err
	}

	if n != 8 {
		return nil, parseEntryErr
	}

	length := binary.BigEndian.Uint64(uIntBuf)

	// decode content
	content := make([]byte, length)

	n, err = buf.Read(content)
	if err != nil {
		return nil, err
	}
	if int(length) != n {
		return nil, parseEntryErr
	}

	// decode verified
	sum := make([]byte, verifiedSize)
	n, err = buf.Read(sum)
	if err != nil {
		return nil, err
	}
	if n != verifiedSize {
		return nil, parseEntryErr
	}

	crc := crc64.New(crc64.MakeTable(verifiedSize))
	crc.Write(content)

	if !tools.EqualSlice(sum, crc.Sum(nil)) {
		return nil, verifiedReqErr
	}

	return &entry{length: length, content: content, ver: sum}, nil
}

const (
	headerSize   = 48
	factor       = 6
	defaultBegin = 256
)

type diskQueue struct {

	// head
	head uint64
	tail uint64

	// The beginning of the first entry
	beginSeek uint64
	length    uint64
	factor    uint64

	// checksum eq (head, tail, beginSeek, length, factor)
	checkSum uint64

	// describe
	// default file is close
	fs *os.File
	mu sync.Mutex

	isClose bool
}

func (d *diskQueue) encodeHeader() []byte {
	enc := bytes.NewBuffer(nil)

	builds := []uint64{d.head, d.tail, d.beginSeek, d.length, d.factor}
	uintBuf := make([]byte, 8)
	for _, ui := range builds {
		binary.BigEndian.PutUint64(uintBuf, ui)
		enc.Write(uintBuf)
	}

	// checkSum codec
	d.checkSum = d.head + d.tail + d.beginSeek + d.length + d.factor
	binary.BigEndian.PutUint64(uintBuf, d.checkSum)
	enc.Write(uintBuf)

	return enc.Bytes()
}

func (d *diskQueue) decodeHeader(buf []byte) error {

	if len(buf) != headerSize {
		return errors.New("non-standard header encoding")
	}

	point := []*uint64{&d.head, &d.tail, &d.beginSeek, &d.length, &d.factor, &d.checkSum}

	for i, p := range point {
		*p = binary.BigEndian.Uint64(buf[i*8 : (i+1)*8])
	}

	if d.checkSum != d.head+d.tail+d.beginSeek+d.length+d.factor {
		return verifiedReqErr
	}

	return nil
}

// 请求逻辑
func (d *diskQueue) Put(req *stream.HttpRequest) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isClose {
		return diskQueueCloseErr
	}

	if d.fs == nil {
		return verifiedReqErr
	}

	ent, err := newEntry(req)
	if err != nil {
		return err
	}

	d.fs.Seek(int64(d.tail), io.SeekStart)
	n, _ := d.fs.Write(encodeDiskEntry(ent))
	d.tail += uint64(n)

	// 重新编码头部
	d.fs.Seek(0, io.SeekStart)
	d.fs.Write(d.encodeHeader())

	d.length += 1
	return nil
}

func (d *diskQueue) Get() (*stream.HttpRequest, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isClose {
		return nil, diskQueueCloseErr
	}

	if d.fs == nil {
		return nil, verifiedReqErr
	}

	if d.length == 0 {
		return nil, errors.New("empty disk queue")
	}

	d.fs.Seek(int64(d.head), io.SeekStart)
	ent, err := getDIskEntry(d.fs)
	if err != nil {
		return nil, err
	}

	d.head += uint64(ent.size())

	// 重新编码头部
	d.fs.Seek(0, io.SeekStart)
	d.fs.Write(d.encodeHeader())

	d.length -= 1

	return stream.DecodeHttpRequest(ent.content)
}

func (d *diskQueue) Size() int {
	return int(d.length)
}

func (d *diskQueue) Cap() int {
	return math.MaxInt
}

func (d *diskQueue) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// 重新编码头部
	d.fs.Seek(0, io.SeekStart)
	d.fs.Write(d.encodeHeader())

	d.isClose = true
	return d.fs.Close()
}

func NewDiskQueue(file string) Buffer {
	disk := new(diskQueue)
	if _, err := os.Stat(file); err != nil {
		disk.head = defaultBegin
		disk.beginSeek = defaultBegin
		disk.tail = disk.head
		disk.length = 0
		disk.factor = factor

		f, err := os.Create(file)
		if err != nil {
			wlog.Fatal(err)
		}

		f.Write(disk.encodeHeader())

		disk.fs = f

		return disk
	}

	head := make([]byte, headerSize)
	f, err := os.OpenFile(file, os.O_RDWR, 0777)
	if err != nil {
		wlog.Fatal(err)
	}

	if n, err := f.Read(head); n != headerSize || err != nil {
		wlog.Fatal("the queue file is corrupted")
	}

	if err := disk.decodeHeader(head); err != nil {
		wlog.Fatal(err)
	}

	disk.fs = f

	return disk
}
