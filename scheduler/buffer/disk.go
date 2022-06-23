package buffer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc64"

	"github.com/wheat-os/slubby/pkg/tools"
	"github.com/wheat-os/slubby/stream"
)

const (
	basicEntrySize  = 8
	entryMinSize    = basicEntrySize * 5
	entryHeaderSize = basicEntrySize * 4

	diskEntryFirstOff = 64 * basicEntrySize
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

type DiskQueue struct {
	firstContentOff int

	head *diskEntry
	tail *diskEntry
}
