package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/wheat-os/slubby/example/temp/spiders"
	"github.com/wheat-os/slubby/example/temp/temp"
	"github.com/wheat-os/wlog"
)

func signalClose(cannel context.CancelFunc) {

	sig := make(chan os.Signal, 1)
	// 监听 退出
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig

	wlog.Info("the listener hears the exit signal, exiting, please do not kill the process directly.")
	cannel()
}

func main() {
	engine := temp.DefaultEngine
	ctx, cannel := context.WithCancel(context.Background())
	go signalClose(cannel)

	engine.Register(spiders.NewTempSpider())

	engine.Start(ctx)

	engine.Close()
}
