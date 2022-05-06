package temp

import (
	"gitee.com/wheat-os/slubby/download"
	"gitee.com/wheat-os/slubby/engine"
	"gitee.com/wheat-os/slubby/outputter"
	"gitee.com/wheat-os/slubby/scheduler"
	"gitee.com/wheat-os/slubby/scheduler/buffer"
	"gitee.com/wheat-os/slubby/scheduler/filter"
)

// Scheduler
var tempFilter = filter.ShortBloomFilter()

var tempBuffer = buffer.ShortQueue()

var tempScheduler = scheduler.ShortScheduler(
	scheduler.WithFilter(tempFilter),
	scheduler.WithBuffer(tempBuffer),
)

// Download
var tempDownload = download.ShortDownload()

// Outputter
var tempOutputter = outputter.ShortOutputter()

// engine
var DefaultEngine = engine.ShortEngine(
	engine.WithScheduler(tempScheduler),
	engine.WithDownload(tempDownload),
	engine.WithOutputter(tempOutputter),
)
