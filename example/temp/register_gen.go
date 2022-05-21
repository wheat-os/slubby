package main

import (
	"github.com/wheat-os/slubby/example/temp/spiders"
	"github.com/wheat-os/slubby/example/temp/temp"
)

func init() {
	engine := temp.DefaultEngine
	engine.Register(spiders.NewTempSpider())
}
