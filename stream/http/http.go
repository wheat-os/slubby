package http

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/wheat-os/slubby/v2/pkg/httputil"
	"github.com/wheat-os/slubby/v2/stream"
)

func init() {
	gob.Register(&stream.Meta{})
}

var (
	ErrNotStreamBody             = errors.New("this http request stream not have body")
	ErrNotExistHttpStreamRequest = errors.New("the input stream is not http stream request")
)

const (
	encStreamHttpRequest = iota
)

type StreamEncoder struct{}

type lowerEnc struct {
	Types int
	Body  []byte
	Ctx   stream.Context
	Cover stream.TargetCover
}

// StreamEncode 实现流编码方法
func (s *StreamEncoder) StreamEncode(stm stream.Stream) ([]byte, error) {
	entry := &lowerEnc{}
	switch v := stm.(type) {
	// stream http request
	case *StreamRequest:
		entry.Cover = v.TargetCover
		entry.Ctx = v.Context
		entry.Types = encStreamHttpRequest
		// body encode
		buf, err := httputil.DumpRequest(v.Request, true)
		if err != nil {
			return nil, err
		}
		entry.Body = buf

	default:
		return nil, stream.ErrUnknownEncoderStreamType
	}

	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(entry); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// StreamDecode 重分离的 binary 中解析流信息
func (s *StreamEncoder) StreamDecode(b []byte) (stream.Stream, error) {
	entry := &lowerEnc{}
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	if err := dec.Decode(entry); err != nil {
		return nil, err
	}

	switch entry.Types {
	// htp request
	case encStreamHttpRequest:
		req, err := httputil.LoadRequest(b)
		if err != nil {
			return nil, err
		}
		return &StreamRequest{TargetCover: entry.Cover, Context: entry.Ctx, Request: req}, nil
	default:
		return nil, stream.ErrUnknownEncoderStreamType
	}
}
