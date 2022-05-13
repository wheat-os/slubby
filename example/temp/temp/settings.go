package temp

import (
	"gitee.com/wheat-os/slubby/download"
	"gitee.com/wheat-os/slubby/download/middle"
	"gitee.com/wheat-os/slubby/engine"
	"gitee.com/wheat-os/slubby/outputter"
	"gitee.com/wheat-os/slubby/outputter/pipline"
	"gitee.com/wheat-os/slubby/scheduler"
	"gitee.com/wheat-os/slubby/scheduler/buffer"
	"gitee.com/wheat-os/slubby/scheduler/filter"
	"gitee.com/wheat-os/wlog"
)

// ***************************************** Logger *****************************************
func init() {
	wlog.SetStdOptions(wlog.WithDisPlayLevel(wlog.InfoLevel))
	wlog.SetStdOptions(wlog.WithDisableCaller(true))
}

// **************************************** Scheduler ***************************************

// filter
var tempFilter = filter.ShortBloomFilter()

// buffer
var tempBuffer = buffer.ShortQueue()

var tempScheduler = scheduler.ShortScheduler(
	scheduler.WithFilter(tempFilter),
	scheduler.WithBuffer(tempBuffer),
)

// **************************************** Download ***************************************
// download middle
var tempMiddleware = middle.MiddleGroup(
	middle.LogMiddle(),
)

var tempDownload = download.ShortDownload(
	download.WithDownloadMiddle(tempMiddleware),
)

// *************************************** Outputter **************************************

// pipline
var tempPipline = pipline.GroupPipline(
	&TempPipline{},
)

var tempOutputter = outputter.ShortOutputter(
	outputter.WithPipline(tempPipline),
)

// ***************************************** Engine ***************************************
var DefaultEngine = engine.ShortEngine(
	engine.WithScheduler(tempScheduler),
	engine.WithDownload(tempDownload),
	engine.WithOutputter(tempOutputter),
)
