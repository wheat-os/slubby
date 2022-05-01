package stream

import "github.com/pkg/errors"

// shortStream
var (
	InvalidEncodingErr = errors.New("decode invalid content, please check your content")
)

// http
var (
	EncodeHttpRequestIsNilErr = errors.New("you cannot encode http request an nil type")
)
