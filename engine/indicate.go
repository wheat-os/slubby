package engine

import "github.com/pkg/errors"

// engine err
var (
	RegisterNilErr = errors.New("registered spiders should not be nil")

	SendUnknownComponentErr = errors.New("the component that needs to be pushed is not recognized")
)

// engine

const (
	StartRequest = "StartRequest should not be given a stream that is not a request"
)
