package tools

import (
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

func TestAntsWlogHandlePanic(t *testing.T) {
	pool, _ := ants.NewPool(10, ants.WithPanicHandler(AntsWlogHandlePanic))

	pool.Submit(func() {
		panic("awdwad")
	})

	time.Sleep(1 * time.Second)
}
