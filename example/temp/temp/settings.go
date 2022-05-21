package temp

import (
	"github.com/wheat-os/slubby/download"
	"github.com/wheat-os/slubby/download/middle"
	"github.com/wheat-os/slubby/engine"
	"github.com/wheat-os/slubby/outputter"
	"github.com/wheat-os/slubby/outputter/pipline"
	"github.com/wheat-os/slubby/scheduler"
	"github.com/wheat-os/slubby/scheduler/buffer"
	"github.com/wheat-os/slubby/scheduler/filter"
	"github.com/wheat-os/wlog"
)

// ***************************************** Logger *****************************************
func init() {
	wlog.SetStdOptions(wlog.WithDisPlayLevel(wlog.InfoLevel))
	wlog.SetStdOptions(wlog.WithDisableCaller(true))
}

// **************************************** Scheduler ***************************************

// filter
var filterModule = filter.ShortBloomFilter()

// buffer
var bufferModule = buffer.ShortQueue()

var schedulerModule = scheduler.ShortScheduler(
	scheduler.WithFilter(filterModule),
	scheduler.WithBuffer(bufferModule),
)

// **************************************** Download ***************************************
// download middle
var downloadMiddlewareModlue = middle.MiddleGroup(
	middle.LogMiddle(),
)

var downloadModule = download.ShortDownload(
	download.WithDownloadMiddle(downloadMiddlewareModlue),
)

// *************************************** Outputter **************************************

// pipline
var piplineModule = pipline.GroupPipline(
	&TempPipline{},
)

var outputterModule = outputter.ShortOutputter(
	outputter.WithPipline(piplineModule),
)

// ***************************************** Engine ***************************************
var DefaultEngine = engine.ShortEngine(
	engine.WithScheduler(schedulerModule),
	engine.WithDownload(downloadModule),
	engine.WithOutputter(outputterModule),
)
