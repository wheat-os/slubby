package error

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestPkgError(t *testing.T) {
	err := errors.New("awd")
	nErr := errors.Wrap(err, "hello")
	fmt.Printf("%v", nErr)
	require.True(t, errors.Is(nErr, err))
}
