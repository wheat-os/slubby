package stream

import (
	"bytes"
	"strings"
	"sync"

	perr "gitee.com/wheat-os/slubby/pkg/error"
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

func SpiderInfo(uid, fqdn string) Stream {
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
		return nil, perr.InvalidEncodingErr
	}

	return &spiderInfo{
		uid:  content[:i],
		fqdn: content[i+len(shortStreamSeparator):],
	}, nil
}

type spiderReflectStream struct {
	refFunc map[string]CallbackFunc
	lock    sync.Mutex
}

func signatureFunc(steam Stream, name string) string {

	const sigSeparator = "."

	uid := "default"
	if steam != nil {
		uid = steam.UId()
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString(uid)
	buf.WriteString(sigSeparator)
	buf.WriteString(name)

	return buf.String()
}

func (s *spiderReflectStream) Register(stream Stream, callbackFunc ...CallbackFunc) error {
	if stream == nil {
		return perr.RegisteredNotSpider
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, callback := range callbackFunc {
		s.refFunc[signatureFunc(stream, callback.Name())] = callback
	}

	return nil
}

// get callback func by steam uid and func name
// used in serialization
// 获取 callback 的反射模型
func (s *spiderReflectStream) CallbackFuncByName(steam Stream, funcName string) CallbackFunc {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.refFunc[signatureFunc(steam, funcName)]

}

func NewReflectStream() *spiderReflectStream {
	return &spiderReflectStream{
		refFunc: make(map[string]CallbackFunc),
	}
}

var std = NewReflectStream()

func MustRegisterSpiderStram(steam Stream, callback ...CallbackFunc) error {
	return std.Register(steam, callback...)
}

func GetCallbackFuncByName(steam Stream, funcName string) CallbackFunc {
	return std.CallbackFuncByName(steam, funcName)
}
