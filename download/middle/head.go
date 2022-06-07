package middle

import "github.com/wheat-os/slubby/stream"

type headMiddle struct {
	hd map[string]string
}

func (h *headMiddle) BeforeDownload(m *M, req *stream.HttpRequest) (*stream.HttpRequest, error) {
	for key, value := range h.hd {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (h *headMiddle) AfterDownload(m *M, req *stream.HttpRequest, resp *stream.HttpResponse) (*stream.HttpResponse, error) {
	return resp, nil
}

func (h *headMiddle) ProcessErr(m *M, req *stream.HttpRequest, err error) {}

func HeadMiddle(head map[string]string) Middleware {
	return &headMiddle{
		hd: head,
	}
}
