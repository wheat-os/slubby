package scheduler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"gitee.com/wheat-os/slubby/stream"
	"gitee.com/wheat-os/wlog"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/require"
)

func putTestRequest(scd Scheduler, count int) error {
	for i := 0; i < count; i++ {
		url := fmt.Sprintf("http://www.test.com/test/%d", i)
		req, err := stream.Request(nil, url, nil)
		if err != nil {
			return err
		}

		if err = scd.Put(req); err != nil {
			return err
		}
	}

	return nil
}

func getTestRequest(scd Scheduler, count int) error {
	for i := 0; i < count; i++ {
		url := fmt.Sprintf("http://www.test.com/test/%d", i)
		req, err := scd.Get()
		if err != nil {
			return err
		}

		if req.URL.String() != url {
			return errors.New("url value err")
		}
	}

	return nil
}

func Test_shortScheduler_Put_Get(t *testing.T) {
	scd := ShortScheduler()

	err := putTestRequest(scd, 100)
	require.NoError(t, err)

	err = getTestRequest(scd, 100)
	require.NoError(t, err)
}

func Test_shortScheduler_Activate(t *testing.T) {
	scd := ShortScheduler()
	require.Equal(t, scd.Activate(), false)
	err := putTestRequest(scd, 1)
	require.NoError(t, err)

	require.Equal(t, scd.Activate(), true)

	err = getTestRequest(scd, 1)
	require.NoError(t, err)

	require.Equal(t, scd.Activate(), false)

}

func Test_shortScheduler_RecvCtxCancel(t *testing.T) {
	cxt, chanel := context.WithCancel(context.Background())

	scd := ShortScheduler()

	ch := scd.RecvCtxCancel(cxt)
	wlog.Info(ch)

	s := scd.(*shortScheduler)
	require.False(t, s.isCancel)

	chanel()
	time.Sleep(time.Second)

	require.True(t, s.isCancel)
	require.False(t, s.Activate())

}

func TestSubmit(t *testing.T) {
	poll, _ := ants.NewPool(5)

	for i := 0; i < 10; i++ {
		poll.Submit(func() {
			fmt.Println(i)
			time.Sleep(time.Second)
		})
	}

	fmt.Println(111111111111111)
}
