package scheduler

import (
	"errors"
	"fmt"
	"testing"

	"gitee.com/wheat-os/slubby/stream"
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
