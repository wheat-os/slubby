package engine

import "github.com/wheat-os/slubby/v2/stream"

type Component interface {
	Streaming(data stream.Stream) error
	Close() error
}
