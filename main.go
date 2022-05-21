package main

import "github.com/wheat-os/wlog"

func main() {
	wlog.SetStdOptions(
		wlog.WithLogLevelColor(wlog.DebugLevel, wlog.FgGreen),
	)

	wlog.Debug("111")
}
