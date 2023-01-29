package scheduler

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/v2/engine"
	"github.com/wheat-os/slubby/v2/stream"
	"testing"
	"time"
)

func TestNewSlubbyScheduler(t *testing.T) {
	tests := []struct {
		name      string
		component engine.SchedulerComponent
	}{
		{
			name: "基本测试",
			component: NewSlubbyScheduler(
				WithCuckooFilterP97("fs.pd"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.component.BackStream()

			const bcKey = "bcKey"

			go func() {
				for {
					stm, cancel := <-ch
					if !cancel {
						return
					}
					fmt.Println(stm.GetMeta(bcKey))
				}
			}()

			for i := 0; i < 100; i++ {
				stm := stream.Background()
				stm.SetMeta(bcKey, i)
				err := tt.component.Streaming(stm)
				require.NoError(t, err)
			}

			time.Sleep(time.Second * 1)

			require.NoError(t, tt.component.Close())
			require.NoError(t, tt.component.Finish())
		})
	}
}
