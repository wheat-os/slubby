package stream

import (
	"bytes"
	"strings"
)

// 控制流
type Stream interface {
	UId() string
	FQDN() string
}

type spiderInfo struct {
	uid  string
	fqdn string
}

func (s *spiderInfo) UId() string {
	return s.uid
}

func (s *spiderInfo) FQDN() string {
	return s.fqdn
}

func shortSpiderInfo(uid, fqdn string) Stream {
	return &spiderInfo{uid: uid, fqdn: fqdn}
}

const (
	shortStreamSeparator = "*$*"
)

// 这个方法编码最低级的 stream
func EncodeShortStream(s Stream) []byte {
	if s == nil || s.UId() == "" {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString(s.UId())
	buf.WriteString(shortStreamSeparator)
	buf.WriteString(s.FQDN())

	return buf.Bytes()
}

func DecodeShortStream(buf []byte) (Stream, error) {
	content := bytes.NewBuffer(buf).String()
	i := strings.Index(content, shortStreamSeparator)
	if i <= 0 {
		return nil, InvalidEncodingErr
	}

	return &spiderInfo{
		uid:  content[:i],
		fqdn: content[i+len(shortStreamSeparator):],
	}, nil
}
