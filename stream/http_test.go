package stream

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/antchfx/htmlquery"
	"github.com/stretchr/testify/require"
)

type TestSpider struct {
	Stream
	id string
}

func (TestSpider) Parse(response *HttpResponse) (Stream, error) {
	fmt.Println("parse")
	return nil, nil
}

func (TestSpider) ToList(response *HttpResponse) (Stream, error) {
	fmt.Println("toList")
	return nil, nil
}

func (TestSpider) GetList(response *HttpResponse) (Stream, error) {
	fmt.Println("getList")
	return nil, nil
}

func TestEncodeCallbackFuncName(t *testing.T) {
	ts := &TestSpider{Stream: shortSpiderInfo("qq", "www.qq.com"), id: "123"}
	req := &HttpRequest{
		Callback: ts.Parse,
	}

	require.Equal(t, req.Callback.Name(), "Parse")
	req.Callback = ts.ToList
	require.Equal(t, req.Callback.Name(), "ToList")
	req.Callback = ts.GetList
	require.Equal(t, req.Callback.Name(), "GetList")

	// 反射回函数
	call := reflect.ValueOf(CallbackFunc(ts.GetList)).Interface()
	// call := reflect.ValueOf(ts.GetList).Interface()
	if callFunc, ok := call.(CallbackFunc); ok {
		callFunc(nil)
	}

}

func equalRequest(t *testing.T, encReq, decReq *HttpRequest) {
	require.Equal(t, decReq.URL, encReq.URL)
	require.Equal(t, decReq.Method, encReq.Method)
	require.Equal(t, decReq.callbackName, encReq.Callback.Name())
	require.Equal(t, decReq.Meta, encReq.Meta)
	// require.Equal(t, decReq.Body, encReq.Body)
	require.Equal(t, decReq.Header, encReq.Header)
}

func TestEncodeHttpRequest(t *testing.T) {
	stearm := shortSpiderInfo("dangdang", "www.baidu.com")
	ts := TestSpider{}
	url := "http://www.baidu.com"

	// Request
	encReq, err := Request(stearm, url, ts.GetList)
	require.NoError(t, err)

	encBuf, err := EncodeHttpRequest(encReq)
	require.NoError(t, err)

	decReq, err := DecodeHttpRequest(encBuf)
	require.NoError(t, err)
	equalRequest(t, encReq, decReq)

	// new Request
	encToBuf, err := EncodeHttpRequest(encReq)
	require.NoError(t, err)

	require.Equal(t, encBuf, encToBuf)

	// BodyRequest
	encReq, err = BodyRequest(stearm, http.MethodPost, url, nil, ts.GetList)
	require.NoError(t, err)
	encBuf, err = EncodeHttpRequest(encReq)
	require.NoError(t, err)
	decReq, err = DecodeHttpRequest(encBuf)
	require.NoError(t, err)
	equalRequest(t, encReq, decReq)

	// FormRequest
	data := map[string][]string{
		"username": {"abb"},
		"password": {"awd"},
	}
	encReq, err = FormRequest(stearm, url, data, ts.Parse)
	require.NoError(t, err)
	encBuf, err = EncodeHttpRequest(encReq)
	require.NoError(t, err)
	decReq, err = DecodeHttpRequest(encBuf)
	require.NoError(t, err)
	equalRequest(t, encReq, decReq)

	// JsonRequest
	jdata := map[string]interface{}{
		"username": "dab",
		"password": "12345644",
	}
	encReq, err = JsonRequest(stearm, url, jdata, ts.GetList)
	require.NoError(t, err)
	encBuf, err = EncodeHttpRequest(encReq)
	require.NoError(t, err)
	decReq, err = DecodeHttpRequest(encBuf)
	require.NoError(t, err)
	equalRequest(t, encReq, decReq)

	// body data and header data request

	// header and body req
	body := `{"data": 123}`
	bodyIo := bytes.NewBufferString(body)
	hReq, err := BodyRequest(nil, http.MethodPost, "http://www.test.com", bodyIo, nil)
	require.NoError(t, err)
	hReq.Header.Add("content-type", "applocation/json")

	// 创建相同的请求
	bodyIo = bytes.NewBufferString(body)
	h2Req, err := BodyRequest(nil, http.MethodPost, "http://www.test.com", bodyIo, nil)
	require.NoError(t, err)
	h2Req.Header.Add("content-type", "applocation/json")

	enc1, err := EncodeHttpRequest(hReq)
	require.NoError(t, err)

	enc2, err := EncodeHttpRequest(h2Req)
	require.NoError(t, err)

	require.Equal(t, enc1, enc2)

}

func TestToBodyEncode(t *testing.T) {
	stearm := shortSpiderInfo("dangdang", "www.baidu.com")
	ts := TestSpider{}
	url := "http://www.baidu.com"

	// JsonRequest
	jdata := map[string]interface{}{
		"username": "dab",
		"password": "12345644",
	}
	encReq, err := JsonRequest(stearm, url, jdata, ts.GetList)
	require.NoError(t, err)
	encBuf, err := EncodeHttpRequest(encReq)
	require.NoError(t, err)
	encToBuf, err := EncodeHttpRequest(encReq)
	require.NoError(t, err)
	require.Equal(t, encBuf, encToBuf)
}

func TestHtmlXpath(t *testing.T) {
	html := `
<div class="nos">
	<p>这几天在构建golang编写的web项目中, 关于dockerfile编写的一些总结</p>
	<a href="https" ima="awd">qq</a> 
</div>
	`

	node, err := htmlquery.Parse(strings.NewReader(html))
	require.NoError(t, err)

	res := htmlquery.Find(node, `//div[@class="nos"]/p/text()`)[0]
	for _, v := range res.Attr {
		fmt.Println(v.Key, v.Val)
	}
	fmt.Println(res.Data)

}

func TestResponseXpath(t *testing.T) {
	html := `
	<div class="nos">
		<p>这几天在构建golang编写的web项目中, 关于dockerfile编写的一些总结</p>
		<a href="https" ima="awd">qq</a> 
	</div>
		`
	res := &HttpResponse{
		Response: &http.Response{},
	}
	res.Body = ioutil.NopCloser(strings.NewReader(html))

	xp := res.Xpath(`//div[@class="nos"]/a)`)[0]
	require.Equal(t, xp.Attribute("href"), "https")
}
