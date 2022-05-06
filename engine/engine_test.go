package engine

import (
	"context"
	"fmt"
	"testing"

	"gitee.com/wheat-os/slubby/stream"
	"gitee.com/wheat-os/wlog"
)

type testSpider struct{}

func (t *testSpider) UId() string {
	return "testID"
}

func (t *testSpider) FQDN() string {
	return "www.baidu.com"
}

func (t *testSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println(response.Text())
	return nil, nil
}

func (t *testSpider) StartRequest() stream.Stream {
	req, err := stream.Request(t, "http://www.baidu.com", nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestShortEngine(t *testing.T) {

	engine := ShortEngine()

	if err := engine.Register(&testSpider{}); err != nil {
		wlog.Fatal(err)
	}

	engine.Start(context.Background())
	engine.Close()

}
