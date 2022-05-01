package spider

import (
	"testing"

	"gitee.com/wheat-os/slubby/stream"
	"github.com/stretchr/testify/require"
)

type testSpider struct {
}

func (t *testSpider) UId() string {
	return "test"
}

func (t *testSpider) FQDN() string {
	return "www.qq.com"
}

func (t *testSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	return stream.Request(t, "http://www.baidu.com", nil)
}

func (t *testSpider) ParseList(response *stream.HttpResponse) (stream.Stream, error) {
	return stream.Request(t, "http://www.qq.com", nil)
}

func (t *testSpider) StartRequest() stream.Stream {
	return nil
}

func TestSpiderManage(t *testing.T) {
	manage := SpiderManage()
	ts := &testSpider{}

	manage.MustRegister(ts)

	resp := &stream.HttpResponse{
		Stream: ts,
	}

	result, err := manage.ParseResp(resp)
	require.NoError(t, err)

	require.Equal(t, result.(*stream.HttpRequest).URL.String(), "http://www.baidu.com")

	manage.RegisterCallbackFunc(ts, ts.Parse)

	result, err = manage.ParseResp(resp)
	require.NoError(t, err)

	require.Equal(t, result.(*stream.HttpRequest).URL.String(), "http://www.baidu.com")

	// 反射模型调用
	resp, err = getParseNameBackFuncResponse(ts, "http://www.qq.com", ts.ParseList)
	require.NoError(t, err)

	result, err = manage.ParseResp(resp)
	require.NoError(t, err)

	require.Equal(t, result.(*stream.HttpRequest).URL.String(), "http://www.qq.com")
}

func getParseNameBackFuncResponse(
	ts Spider,
	url string,
	callback stream.CallbackFunc,
) (*stream.HttpResponse, error) {
	// 反射模型调用
	req, err := stream.Request(ts, url, callback)
	if err != nil {
		return nil, err
	}

	// 通过 编解码，使 resp 带有 parseName 参数
	buf, err := stream.EncodeHttpRequest(req)
	if err != nil {
		return nil, err
	}
	req, err = stream.DecodeHttpRequest(buf)
	if err != nil {
		return nil, err
	}

	resp := &stream.HttpResponse{
		Stream: ts,
	}

	resp.WithHttpAndRequestStream(req, nil)

	return resp, nil
}

func BenchmarkRegisterBackFunc(b *testing.B) {
	b.ResetTimer()
	manage := SpiderManage()
	ts := &testSpider{}

	manage.MustRegister(ts)
	manage.RegisterCallbackFunc(ts, ts.ParseList)

	resp, err := getParseNameBackFuncResponse(ts, "http://www.baidu.com", ts.ParseList)
	require.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		manage.ParseResp(resp)
		require.NoError(b, err)
	}

}

func BenchmarkRefBackFunc(b *testing.B) {
	b.ResetTimer()
	manage := SpiderManage()
	ts := &testSpider{}

	manage.MustRegister(ts)
	// manage.RegisterCallbackFunc(ts, ts.ParseList)
	resp, err := getParseNameBackFuncResponse(ts, "http://www.baidu.com", ts.ParseList)
	require.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		manage.ParseResp(resp)
		require.NoError(b, err)
	}
}
