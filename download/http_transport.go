package download

import (
	"github.com/wheat-os/slubby/v2/stream"
	sHttp "github.com/wheat-os/slubby/v2/stream/http"
	"net/http"
)

func withHttpClientTransport(client *http.Client) OptFunc {
	if client == nil {
		client = http.DefaultClient
	}

	roundTripper := func(inStream stream.Stream) (stream.Stream, error) {
		req, ok := inStream.(*sHttp.StreamRequest)
		if !ok {
			return nil, sHttp.ErrNotExistHttpStreamRequest
		}
		resp, err := client.Do(req.Request)
		if err != nil {
			return nil, err
		}
		return sHttp.NewHttpResponse(resp)
	}

	return func(opt *SlubbyComponent) {
		opt.roundTripper = roundTripper
	}
}

// WithHttpTransport 使用 go http transport
func WithHttpTransport(transport http.RoundTripper) OptFunc {
	if transport == nil {
		transport = http.DefaultTransport
	}

	client := &http.Client{Transport: transport}
	return withHttpClientTransport(client)
}

// WithHttpClientTransport 使用 go http client
func WithHttpClientTransport(client *http.Client) OptFunc {
	return withHttpClientTransport(client)
}
