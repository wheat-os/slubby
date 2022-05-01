package spider

import "github.com/pkg/errors"

var (
	RegisteredNotSpider       = errors.New("the registered type must be SpiderInfo")
	RegisterSpiderUidConflict = errors.New("uid already exists in the crawler collection, try replacing spider uid")

	RefCallFuncBackErr = errors.New("an unreliable reflection function was called")
)
