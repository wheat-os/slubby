package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testStream struct{}

func (t *testStream) UId() string {
	panic("not implemented") // TODO: Implement
}

func (t *testStream) FQDN() string {
	return "test"
}

func TestShortLimiter(t *testing.T) {
	lim := ShortLimiter(time.Second * 1)

	require.Equal(t, lim.Allow(&testStream{}), true)

	// 1 second
	require.Equal(t, lim.Allow(&testStream{}), true)

}
