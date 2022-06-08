package pipline

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/wheat-os/slubby/stream"
)

type csvPIpline struct {
	path      string
	frequency int
	counter   int
	isCreate  bool
	f         *os.File
	headers   []string
	buf       *bytes.Buffer
	mu        sync.Mutex
}

func (c *csvPIpline) fash() error {
	_, err := io.Copy(c.f, c.buf)
	return err
}

func (c *csvPIpline) write(content []string) {
	c.buf.WriteString(strings.Join(content, ","))
	c.buf.WriteString("\n")
}

func (c *csvPIpline) OpenSpider() error {
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		c.isCreate = true
	}

	fs, err := os.OpenFile(c.path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	c.f = fs
	c.buf = bytes.NewBuffer(nil)

	return nil
}

func (c *csvPIpline) CloseSpider() error {
	defer c.f.Close()
	return c.fash()
}

func (c *csvPIpline) ProcessItem(item stream.Item) stream.Item {
	typeOfCat := reflect.TypeOf(item)

	if typeOfCat.Kind() == reflect.Ptr {
		typeOfCat = typeOfCat.Elem()
	}

	valueOfCat := reflect.ValueOf(item)

	if valueOfCat.Kind() == reflect.Ptr {
		valueOfCat = valueOfCat.Elem()
	}

	headers := make([]string, 0, 10)
	bodes := make([]string, 0, 10)
	for i := 0; i < typeOfCat.NumField(); i++ {
		field := typeOfCat.Field(i)
		if fileName := field.Tag.Get("csv"); fileName != "" {
			headers = append(headers, fileName)
			bodes = append(bodes, fmt.Sprintf("%v", valueOfCat.Field(i).Interface()))
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 头部数据
	if len(c.headers) == 0 && c.isCreate {
		c.headers = headers
		c.write(headers)
		c.fash()
	}

	c.write(bodes)
	if c.counter%c.frequency == 0 {
		c.fash()
	}

	return nil
}

func SaveCsvPipline(path string, freq int) Pipline {
	if freq <= 1 {
		freq = 1
	}
	return &csvPIpline{
		path:      path,
		frequency: freq,
		counter:   0,
	}
}
