package spiders

import (
	"fmt"

	"gitee.com/wheat-os/slubby/stream"
)

type TempSpider struct{}

func (t *TempSpider) UId() string {
	return "temp"
}

func (t *TempSpider) FQDN() string {
	return "www.temp.com"
}

func (t *TempSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println(response.Status)

	return nil, nil
}

func (t *TempSpider) StartRequest() stream.Stream {
	req, _ := stream.Request(t, "http://www.baidu.com", nil)
	return req
}
