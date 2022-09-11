package stream

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	urls "net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/wheat-os/wlog"
	"golang.org/x/net/html"
)

type CallbackFunc func(response *HttpResponse) (Stream, error)

// func name
func (c CallbackFunc) Name() string {
	// 获取函数名
	fn := runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name()
	nameGex := regexp.MustCompile(`.+\.(\w+)-fm`)
	names := nameGex.FindAllStringSubmatch(fn, -1)

	if len(names) == 1 && len(names[0]) == 2 {
		return names[0][1]
	}

	return fn
}

type HttpRequest struct {
	*http.Request

	// 爬虫回调函数
	Callback CallbackFunc
	// 跳过过滤器
	SkipFilter bool
	// 带上请求的一些参数
	Meta map[string]interface{}
	// 下载失败重新尝试次数
	Retry int

	// 爬虫流
	stream Stream

	// 反射后调用的回调函数名称
	callbackName string
}

// 几种 http new func
func Request(self Stream, url string, callback CallbackFunc) (*HttpRequest, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return &HttpRequest{Request: req, Callback: callback, stream: self, Meta: make(map[string]interface{})}, nil
}

func BodyRequest(
	self Stream,
	method string,
	url string,
	body io.Reader,
	callback CallbackFunc,
) (*HttpRequest, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return &HttpRequest{Request: req, Callback: callback, stream: self, Meta: make(map[string]interface{})}, nil
}

func FormRequest(
	self Stream,
	url string,
	fromData map[string][]string,
	callback CallbackFunc,
) (*HttpRequest, error) {
	data := urls.Values(fromData)
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return &HttpRequest{Request: req, stream: self, Meta: map[string]interface{}{}, Callback: callback}, nil
}

func JsonRequest(
	self Stream,
	url string,
	jsonData map[string]interface{},
	callback CallbackFunc,
) (*HttpRequest, error) {
	data, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	dataBuf := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, url, dataBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	return &HttpRequest{Request: req, stream: self, Meta: make(map[string]interface{}), Callback: callback}, nil
}

func (h *HttpRequest) UId() string {
	if h.stream != nil {
		return h.stream.UId()
	}

	return ""
}

func (h *HttpRequest) FQDN() string {
	if h.stream != nil {
		return h.stream.FQDN()
	}

	return ""
}

// 用于持久化的 request
type shortRequest struct {
	Url          string
	Method       string
	Header       map[string][]string
	Meta         map[string]interface{}
	Body         []byte
	SpiderInfo   []byte
	CallBackName string
}

// 降低请求等级
func (s *shortRequest) buildShortRequest(req *HttpRequest) error {
	if req == nil {
		return EncodeHttpRequestIsNilErr
	}

	s.Url = req.URL.String()
	s.Method = req.Method
	s.Header = req.Header

	// 构建 shortRequest
	body := bytes.NewBuffer(nil)
	if req.Body != nil {
		if _, err := io.Copy(body, req.Body); err != nil {
			return err
		}
		s.Body = body.Bytes()

		// 重置 body
		req.Body = ioutil.NopCloser(body)
	}

	// 解析 callback func方法
	if req.Callback != nil {
		s.CallBackName = req.Callback.Name()
	}

	// 构建 spider info
	s.SpiderInfo = EncodeShortStream(req.stream)

	return nil
}

// 升级请求
func (s *shortRequest) ToHttpRequest() (*HttpRequest, error) {
	stream, err := DecodeShortStream(s.SpiderInfo)

	// default stream, stream is nil
	if err == InvalidEncodingErr {
		stream = shortSpiderInfo("", "")
	}

	var body io.Reader = nil
	if len(s.Body) != 0 {
		body = bytes.NewBuffer(s.Body)
	}

	req, err := BodyRequest(stream, s.Method, s.Url, body, nil)
	if err != nil {
		return nil, err
	}
	req.callbackName = s.CallBackName
	req.Header = s.Header
	req.Meta = s.Meta

	return req, nil
}

func EncodeHttpRequest(req *HttpRequest) ([]byte, error) {
	shortReq := new(shortRequest)
	if err := shortReq.buildShortRequest(req); err != nil {
		return nil, err
	}

	// gob encode shortHttpRequest
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(shortReq); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeHttpRequest(b []byte) (*HttpRequest, error) {
	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	eReq := new(shortRequest)
	err := dec.Decode(eReq)
	if err != nil {
		return nil, err
	}

	return eReq.ToHttpRequest()
}

type HttpResponse struct {
	*http.Response
	// 带上请求的一些参数
	Meta map[string]interface{}

	Stream

	// 解析函数
	Callback CallbackFunc

	// 通过反射方案执行
	parseFuncName string
}

func (h *HttpResponse) WithHttpAndRequestStream(sReq *HttpRequest, req *http.Response) {
	h.Response = req
	h.Meta = sReq.Meta
	h.Callback = sReq.Callback
	h.Stream = sReq.stream
	h.parseFuncName = sReq.callbackName
}

func (h *HttpResponse) ParseName() string {
	return h.parseFuncName
}

func (h *HttpResponse) Bytes() []byte {
	buf, _ := ioutil.ReadAll(h.Body)
	return buf
}

func (h *HttpResponse) Text() string {
	buf, _ := ioutil.ReadAll(h.Body)
	return string(buf)
}

func (h *HttpResponse) Xpath(expr string) []*xpath {
	node, err := htmlquery.Parse(h.Body)
	if err != nil {
		wlog.Error(err)
		return nil
	}

	return (*xpath)(node).Xpath(expr)
}

type xpath html.Node

func (x *xpath) Xpath(expr string) []*xpath {
	nodes := htmlquery.Find((*html.Node)(x), expr)
	xs := make([]*xpath, 0, len(nodes))
	for _, v := range nodes {
		xs = append(xs, (*xpath)(v))
	}

	return xs
}

func (x *xpath) Attribute(key string) string {
	for _, attr := range x.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}
