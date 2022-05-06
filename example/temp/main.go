package main

import (
	"context"

	"gitee.com/wheat-os/slubby/example/temp/spiders"
	"gitee.com/wheat-os/slubby/example/temp/temp"
)

func main() {
	engine := temp.DefaultEngine

	engine.Register(&spiders.TempSpider{})

	engine.Start(context.Background())

	engine.Close()
}
