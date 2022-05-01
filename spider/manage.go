package spider

import (
	"fmt"
	"reflect"
	"sync"

	"gitee.com/wheat-os/slubby/stream"
	"github.com/pkg/errors"
)

type reflectSpider struct {
	self    reflect.Value
	memFunc map[string]stream.CallbackFunc
	refFunc map[string]reflect.Value
}

type spiderManage struct {
	manage map[string]*reflectSpider
	lock   sync.RWMutex
}

// register Spider
func (s *spiderManage) MustRegister(sp Spider) {

	s.lock.Lock()
	defer s.lock.Unlock()

	// TODO repeat uid err

	if sp == nil {
		return
	}

	spVal := reflect.ValueOf(sp)
	spType := spVal.Type()

	refFunc := make(map[string]reflect.Value)

	for i := 0; i < spVal.NumMethod(); i++ {
		if parseName := spType.Method(i).Name; parseName != "" {
			refFunc[parseName] = spType.Method(i).Func
		}
	}

	spiderRef := &reflectSpider{
		self:    spVal,
		memFunc: make(map[string]stream.CallbackFunc),
		refFunc: refFunc,
	}

	s.manage[sp.UId()] = spiderRef
}

func (s *spiderManage) RegisterCallbackFunc(sp Spider, callback stream.CallbackFunc) {
	s.lock.Lock()
	defer s.lock.Unlock()
	refSpider := s.manage[sp.UId()]
	if refSpider == nil {
		s.lock.Unlock()
		s.MustRegister(sp)
		s.lock.Lock()
	}

	if refSpider != nil {
		refSpider.memFunc[callback.Name()] = callback
	}
}

// 爬虫解析方案
func (s *spiderManage) ParseResp(resp *stream.HttpResponse) (stream.Stream, error) {
	// 优先执行 callback
	if resp.Callback != nil {
		return resp.Callback(resp)
	}

	// 尝试使用反射执行
	s.lock.RLock()
	defer s.lock.RUnlock()
	if funcName := resp.ParseName(); funcName != "" {
		if refFunc := s.manage[resp.UId()]; refFunc != nil {
			// 尝试不使用反射方案
			if callback := refFunc.memFunc[funcName]; callback != nil {
				return callback(resp)
			}

			if refBack := refFunc.refFunc[funcName]; !refBack.IsNil() {
				result := refBack.Call([]reflect.Value{refFunc.self, reflect.ValueOf(resp)})

				resErr := func() error {
					msg := fmt.Sprintf(
						"callback func uid:%s, func name: %s",
						resp.UId(),
						resp.ParseName(),
					)

					return errors.Wrap(RefCallFuncBackErr, msg)
				}
				if len(result) != 2 {
					return nil, resErr()
				}

				resultStream, ok := result[0].Interface().(stream.Stream)
				if !ok {
					return nil, resErr()
				}

				if result[1].IsNil() {
					return resultStream, nil
				}

				resultErr, ok := result[1].Interface().(error)
				if !ok {
					return nil, resErr()
				}

				return resultStream, resultErr
			}
		}
	}

	// 默认执行
	if info, ok := resp.Stream.(Spider); ok {
		return info.Parse(resp)
	}

	return nil, nil
}

func SpiderManage() *spiderManage {
	return &spiderManage{
		manage: make(map[string]*reflectSpider),
	}
}
