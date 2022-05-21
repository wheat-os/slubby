package spiders

import (
	"fmt"
	"sync"

	"gitee.com/wheat-os/slubby/spider"
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
	fmt.Println(response.Text())

	return nil, nil
}

func (t *TempSpider) StartRequest() stream.Stream {
	req, _ := stream.Request(t, "http://www.baidu.com", nil)
	return req
}

var (
	once = sync.Once{}
	temp *TempSpider
)

func NewTempSpider() spider.Spider {
	once.Do(func() {
		temp = &TempSpider{}
	})
	return temp
}
