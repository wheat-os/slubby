package middle

import (
	"github.com/wheat-os/slubby/stream"
	"github.com/wheat-os/wlog"
)

type logMiddle struct{}

func (l *logMiddle) BeforeDownload(m *M, req *stream.HttpRequest) (*stream.HttpRequest, error) {
	return req, nil
}

func (l *logMiddle) AfterDownload(m *M, req *stream.HttpRequest, resp *stream.HttpResponse) (*stream.HttpResponse, error) {

	wlog.Infof("<Response [%d]> Request<url: %s, method: %s",
		resp.StatusCode, req.URL.String(), req.Method)

	return resp, nil
}

func (l *logMiddle) ProcessErr(m *M, req *stream.HttpRequest, err error) {
	wlog.Errorf("<Request url: %s> download err, err: %s", req.URL.String(), err)
}

func LogMiddle() Middleware {
	return &logMiddle{}
}
