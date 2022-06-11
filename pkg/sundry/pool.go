package sundry

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/wheat-os/wlog"
)

func AntsWlogHandlePanic(err interface{}) {
	if err == nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s\n", err))

	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buf.WriteString(fmt.Sprintf("%s:%d\n", file, line))
	}

	wlog.Error(buf.String())
}
