package buffer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc64"
	"io"
	"math"
	"os"
	"sync"

	"github.com/wheat-os/slubby/pkg/tools"
	"github.com/wheat-os/slubby/stream"
)

const (
	basicEntrySize  = 8
	entryMinSize    = basicEntrySize * 5
	entryHeaderSize = basicEntrySize * 4
)

var (
	errEntryIllegal = errors.New("an unhandled error occurred to parse an entry")
)

type diskEntry struct {

	// entry offset
	start uint64
	end   uint64
	next  uint64

	// length content
	length  uint64
	content []byte
	ver     []byte
}

func encodeDiskEntry(ent *diskEntry) []byte {
	uintByte := make([]byte, basicEntrySize)
	buf := bytes.NewBuffer(nil)
	crc := crc64.New(crc64.MakeTable(basicEntrySize))

	for _, val := range []uint64{ent.start, ent.end, ent.next, ent.length} {
		binary.BigEndian.PutUint64(uintByte, val)
		buf.Write(uintByte)
		crc.Write(uintByte)
	}

	// content
	buf.Write(ent.content)
	crc.Write(ent.content)

	ent.ver = crc.Sum(nil)
	buf.Write(ent.ver)

	return buf.Bytes()
}

func decodeDiskEntry(b []byte) (*diskEntry, error) {
	if len(b) < entryMinSize {
		return nil, errEntryIllegal
	}

	crc := crc64.New(crc64.MakeTable(basicEntrySize))
	ent := new(diskEntry)

	for i, p := range []*uint64{&ent.start, &ent.end, &ent.next, &ent.length} {
		uintByte := b[basicEntrySize*i : (i+1)*basicEntrySize]
		*p = binary.BigEndian.Uint64(uintByte)
		crc.Write(uintByte)
	}

	// check entry value
	if int(ent.end-ent.start) != len(b) {
		return nil, errEntryIllegal
	}

	if len(b)-entryMinSize != int(ent.length) {
		return nil, errEntryIllegal
	}

	// ver
	content := b[entryHeaderSize : entryHeaderSize+ent.length]
	sum := b[entryHeaderSize+ent.length:]
	crc.Write(content)
	if !tools.EqualSlice(crc.Sum(nil), sum) {
		return nil, errEntryIllegal
	}

	ent.content = content
	ent.ver = sum
	return ent, nil
}

func newDiskEntryByHttpRequest(req *stream.HttpRequest, start uint64) (*diskEntry, error) {
	reqBuf, err := stream.EncodeHttpRequest(req)
	if err != nil {
		return nil, err
	}

	length := uint64(len(reqBuf))
	end := start + entryMinSize + length

	return &diskEntry{start: start, end: end, length: length, content: reqBuf}, nil
}

const (
	diskEntryFirstOff = 64 * basicEntrySize
	headOff           = 3 * basicEntrySize
)

var (
	errSetEntryToDisk = errors.New("an error occurred in entering the disk")
	errDiskQueueEmpty = errors.New("the current disk queue is empty")
	errParseDiskEntry = errors.New("parses the entry error")
)

type diskQueue struct {
	firstContentOff uint64
	length          uint64

	head *diskEntry
	tail *diskEntry

	ver []byte

	mu sync.Mutex

	fs *os.File
}

func (d *diskQueue) updateEntryNext(ent *diskEntry) {
	d.fs.Seek(int64(ent.start+2*basicEntrySize), io.SeekStart)

	uintByte := make([]byte, basicEntrySize)
	binary.BigEndian.PutUint64(uintByte, ent.next)
	d.fs.Write(uintByte)
}

func (d *diskQueue) setEntry(ent *diskEntry) error {
	d.fs.Seek(int64(ent.start), io.SeekStart)
	buf := encodeDiskEntry(ent)
	if n, _ := d.fs.Write(buf); n != int(ent.end-ent.start) {
		return errEntryIllegal
	}

	return nil
}

func (d *diskQueue) readEntryByOff(start uint64) (*diskEntry, error) {
	d.fs.Seek(int64(start), io.SeekStart)
	head := make([]byte, entryHeaderSize)
	entStart := binary.BigEndian.Uint64(head[0*basicEntrySize : basicEntrySize*1])
	entEnd := binary.BigEndian.Uint64(head[1*basicEntrySize : basicEntrySize*2])

	if entStart != start {
		return nil, errParseDiskEntry
	}

	d.length -= 1

	body := make([]byte, entEnd-entStart)
	d.fs.Read(body)

	return decodeDiskEntry(append(head, body...))
}

func (d *diskQueue) flashHeader() {
	buf := bytes.NewBuffer(nil)
	uintByte := make([]byte, basicEntrySize)

	// firstContentOff
	binary.BigEndian.PutUint64(uintByte, d.firstContentOff)
	buf.Write(uintByte)

	// length
	binary.BigEndian.PutUint64(uintByte, d.length)
	buf.Write(uintByte)

	buf.Write(encodeDiskEntry(d.head))
	buf.Write(encodeDiskEntry(d.tail))

	crc := crc64.New(crc64.MakeTable(basicEntrySize))
	crc.Write(buf.Bytes())

}

// 请求逻辑
func (d *diskQueue) Put(req *stream.HttpRequest) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	ent, err := newDiskEntryByHttpRequest(req, d.tail.end)
	if err != nil {
		return err
	}

	if err = d.setEntry(ent); err != nil {
		return err
	}

	// set tail
	d.tail.next = d.tail.end
	d.updateEntryNext(d.tail)

	d.length += 1
	d.tail = ent

	return nil
}

func (d *diskQueue) Get() (*stream.HttpRequest, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.length <= 0 || d.head.next == 0 {
		return nil, errDiskQueueEmpty
	}

	ent, err := d.readEntryByOff(d.head.next)
	if err != nil {
		return nil, err
	}

	d.head = ent
	d.length -= 1
	return stream.DecodeHttpRequest(ent.content)
}

func (d *diskQueue) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return int(d.length)
}

func (d *diskQueue) Cap() int {
	return math.MaxInt
}

func (d *diskQueue) Close() error {
	panic("not implemented") // TODO: Implement
}

func DiskQueue(file string) Buffer {
	// append
	if _, err := os.Stat(file); err != nil {
		fs, err := os.Create(file)
		tools.WlogPanicErr(err)
		headEntry := &diskEntry{
			start:   diskEntryFirstOff,
			end:     diskEntryFirstOff,
			next:    diskEntryFirstOff,
			length:  0,
			content: nil,
			ver:     nil,
		}

		tailEntry := &diskEntry{
			start:   diskEntryFirstOff,
			end:     diskEntryFirstOff,
			next:    diskEntryFirstOff,
			length:  0,
			content: nil,
			ver:     nil,
		}

		return &diskQueue{
			firstContentOff: diskEntryFirstOff,
			length:          0,
			fs:              fs,
			head:            headEntry,
			tail:            tailEntry,
		}
	}

	// fs, err := os.OpenFile(file, os.O_RDWR, 0777)
	// tools.WlogPanicErr(err)

	// // read header
	return nil
}
