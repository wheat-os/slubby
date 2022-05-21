package spider

import "github.com/wheat-os/slubby/stream"

type Spider interface {
	stream.Stream

	Parse(response *stream.HttpResponse) (stream.Stream, error)
	StartRequest() stream.Stream
}
