package main

import (
	"gitee.com/wheat-os/slubby/example/temp/spiders"
	"gitee.com/wheat-os/slubby/example/temp/temp"
)

func init() {
	engine := temp.DefaultEngine
	engine.Register(spiders.NewTempSpider())
}
