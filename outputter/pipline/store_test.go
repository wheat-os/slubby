package pipline

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wheat-os/slubby/stream"
)

type storeItem struct {
	stream.Item
	Name    string `csv:"name" json:"name"`
	Content string `csv:"content"`
	Abb     int    `csv:"abb"`
}

// process pipline
func Test_csvPIpline_ProcessItem(t *testing.T) {
	testName := "wad"
	testContent := "awd"
	item := &storeItem{
		Item:    stream.BasicItem(nil),
		Name:    testName,
		Content: testContent,
	}

	path := "./test.csv"

	pip := SaveCsvPipline(path, 20)
	pip.OpenSpider()
	pip.ProcessItem(item)
	pip.ProcessItem(item)
	pip.CloseSpider()

	result := `name,content,abb
wad,awd,0
wad,awd,0
`

	res, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, result, string(res))

	os.Remove(path)
}
