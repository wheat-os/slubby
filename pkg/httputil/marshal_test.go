package httputil

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func buildHttpRequest(t *testing.T) *http.Request {
	body := bytes.NewBuffer([]byte("awhjidawd"))
	req, err := http.NewRequest(http.MethodGet, "www.baidu.com", body)
	require.NoError(t, err)
	//req.Header.Set("use", "ybb3")
	//req.Header.Set("use121", "ybb3")
	//req.Header.Set("usjklhawde", "ybb3")
	//req.Header.Set("useiawd&bawd", "ybb3aiowd")
	return req
}

func TestDumpRequest(t *testing.T) {
	req := buildHttpRequest(t)
	buf, err := DumpRequest(req, true)
	require.NoError(t, err)
	fmt.Println(string(buf))
}

func TestLoadRequest(t *testing.T) {
	req := buildHttpRequest(t)
	buf, err := DumpRequest(req, true)
	require.NoError(t, err)
	req, err = LoadRequest(buf)
	require.NoError(t, err)
	buf2, err := DumpRequest(req, true)
	require.NoError(t, err)
	require.Equal(t, buf, buf2)
}
