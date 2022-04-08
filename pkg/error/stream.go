package error

import "github.com/pkg/errors"

// shortStream
var (
	InvalidEncodingErr        = errors.New("decode invalid content, please check your content")
	RegisteredNotSpider       = errors.New("the registered type must be SpiderInfo")
	RegisterSpiderUidConflict = errors.New("uid already exists in the crawler collection, try replacing spider uid")
)

// http
var (
	EncodeHttpRequestIsNilErr = errors.New("you cannot encode http request an nil type")
)
