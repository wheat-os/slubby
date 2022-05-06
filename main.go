package main

import "gitee.com/wheat-os/wlog"

func main() {
	wlog.SetStdOptions(
		wlog.WithLogLevelColor(wlog.DebugLevel, wlog.FgGreen),
	)

	wlog.Debug("111")
}
