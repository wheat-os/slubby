package spider

import (
	"fmt"

	"gitee.com/wheat-os/slubby/stream"
)

type TestSpider struct {
	stream.Stream
	id string
}

func (TestSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println("parse")
	return nil, nil
}

func (TestSpider) ToList(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println("toList")
	return nil, nil
}

func (TestSpider) GetList(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println("getList")
	return nil, nil
}
