# slubby 使用文档

这个文档会介绍 slubby 以及配套脚手架框架 slub 的使用方法，部分特性以及更多高级的组件扩展方案，帮助使用者了解整个 slubby 框架。

### 安装

本部分我们主要来安装使用 slub 手脚架（一个用来快速开发 slubby 项目的应用框架），在安装框架以前，我们需要确定是否拥有支持环境，以下环境是必须的。

- go 1.15+ （可以使用 go vesion 检查版本）
  
- go mod
  

确定环境满足以后我们可以使用以下命令来安装 安装 最新版本的 slub

```shell
go install github.com/wheat-os/slub@latest
```

或者我们使用源代码编译安装
```sh
# 如果无法访问 github 可以用 gitee.com/wheat-os/slub
git clone github.com/wheat-os/slub
go build

# windows 可以配置环境变量
cp ./slub /usr/bin
```

国内开发者可能出现网络错误，可以设置 go mod 代理，[go proxy](https://goproxy.cn/) 。安装完成以后我们可以执行 `slub` 出现 slub 介绍页面表示 slub 工具安装成功。好的，我想你一定成功了对吧，这非常简单，我们现在进入下一步。

### 快速体验

##### 创建项目

我们可以使用 `slub startproject <name>` 来创建一个 slubby 项目，如 `slub startproject first` 我们会得到这样一个目录结构

```shell
├── conf.toml
├── first
│   ├── items.go
│   ├── middleware.go
│   ├── piplines.go
│   └── settings.go
├── go.mod
├── main.go
└── spiders
```

目录中每个文件的作用，我们会在以后继续学习，现在，我们开始创建第一个爬虫实例，来开始我们的 slubby 旅行。

现在我们没有安装项目需要的 package 支持，我们可以在项目目录中执行 `go mod tidy` 来安装项目需要的包支持（go1.17 执行失败可以尝试 `go mod tidy -compat=1.17`），现在相信你一定成功的解决了项目中的包错误，我们进入下一步。

##### 创建爬虫

我们使用 `slub genspider <spider name> <fqdn>` 来创建一个爬虫模板，如 `slub genspider first www.baidu.com` ，注意我们必须进入爬虫项目目录来执行命令。现在我们检查 spiders 目录中的 `first.go` 文件。

```go
package spiders

import (
	"fmt"
	"sync"

	"github.com/wheat-os/slubby/spider"
	"github.com/wheat-os/slubby/stream"
)

type firstSpider struct{}

// uid 是每一个爬虫的唯一识别码，我们应该保证同一个 slubby 项目中它是唯一的。
func (t *firstSpider) UId() string {
	return "first"
}

// fqdn 记录了这个爬虫工作的域名，用于帮助下载器限制器工作。
func (t *firstSpider) FQDN() string {
	return "www.baidu.com"
}

// Parse 是默认的爬虫解析函数，当一个请求没有指定回调函数时，默认调用它。
func (t *firstSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	fmt.Println(response.Text())
    
	return nil, nil
}

// StartRequest 是爬虫的入口函数，我们会在这里返回爬虫的开始请求。
func (t *firstSpider) StartRequest() stream.Stream {
	req, _ := stream.Request(t, "http://www.baidu.com", nil)
	return req
}

var (
	firstOnce = sync.Once{}
	first     *firstSpider
)

// 单例的爬虫创建函数
func FirstSpider() spider.Spider {
	firstOnce.Do(func() {
		first = &firstSpider{}
	})
	return first
}

```

我们现在对爬虫进行修改来获取 百度的 `title: 百度一下，你就知道`，我们主要修改 `Parse` 来解析 百度给的返回值。

```go
func (t *firstSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	content := response.Text()

	reg, _ := regexp.Compile(`<title>(.*?)</title>`)

	result := reg.FindAllStringSubmatch(content, 1)
	wlog.Debug(result[0][1])

	return nil, nil
}
```

这里我们使用正则提取了 response 中 的 title 信息，并且通过 `wlog(slub 使用的日志工具)` 打印到控制台。

现在我们来运行这个爬虫，我们会介绍 2 个运行方法来运行我们的爬虫实例。

##### 运行爬虫

**使用 go 编译运行**

我们可以执行 `go run .` 来运行项目，这样我们不会看到任何的输出，这是由于我们并没有向项目注册爬虫，slubby 框架无法解析运行我们定义的 `first spider`，我们需要再 `main.go` 中去注册我们定义的 spider。

```go
func main() {
	viper.SetConfigFile(projectConfFile)

	engine := first.DefaultEngine

	// 注册 spider
	engine.Register(spiders.FirstSpider())

	ctx, cannel := context.WithCancel(context.Background())
	go signalClose(cannel)

	engine.Start(ctx)

	engine.Close()
}

```

我们再次执行 `go run .`

```
2022-05-27 11:00:40 INFO <Response [200]> Request<url: http://www.baidu.com, method: GET
2022-05-27 11:00:40 DEBUG 百度一下，你就知道
2022-05-27 11:00:42 INFO the spider shuts down successfully
```

天啊，这样岂不是我们每次创建爬虫都需要去注册它，有更好的办法吗，那么现在我们删除掉注册爬虫的代码 `engine.Register(spiders.FirstSpider())` 尝试使用 `slub` 工具来运行爬虫。

**使用 slub 运行**

我们可以再项目中执行 `slub crawl` 来运行项目中全部的爬虫，如果我们只希望运行部分爬虫可以使用 `slub crawl --run=first,second` 的方式来运行 first 以及 second 两个爬虫。如我们执行 `slub crawl --run=first` 或者是 `slub crawl`

```
2022-05-27 11:06:31 INFO <Response [200]> Request<url: http://www.baidu.com, method: GET
2022-05-27 11:06:31 DEBUG 百度一下，你就知道
2022-05-27 11:06:33 INFO the spider shuts down successfully
```

得到一样的结果。

**这里发生了什么**

slubby 会把组成到它上的爬虫执行 `StartRequest` 来获取 `请求流` 并且将它发送给下载器，由下载器下载完成后，把 response 交给爬虫解析。

上述的例子中我们直接再爬虫解析函数中打印了目标数据，但是如果我们需要做更多的 IO 操作来保存数据 slubby 并不建议再爬虫解析部分调用 IO 操作。这时候可以使用到 slubby 提供的输出器。

##### 输出器

再了解输出器以前，我们先定义一个 数据流（`item`) slub 项目中的 `item` 统一在 项目目录里的 `item.go` 中管理。我们添加一个数据字段 title (注意由于 go 的特性应该大写)。

```go
package first

import "github.com/wheat-os/slubby/stream"

type FirstItem struct {
	stream.Item
	Title string
}
 
```

现在我们优化 `Parse` 函数的写法，不在 spider 中执行 IO 操作。

我们更新 spider `Parse` 函数返回一个 `item`

```go
func (t *firstSpider) Parse(response *stream.HttpResponse) (stream.Stream, error) {
	content := response.Text()

	reg, _ := regexp.Compile(`<title>(.*?)</title>`)

	result := reg.FindAllStringSubmatch(content, 1)

	item := &first.FirstItem{
		Item:  stream.BasicItem(t),
		Title: result[0][1],
	}

	return item, nil
}
```

上述代码中 `Item: stream.BasicItem` 是必须，在 Parse 中返回的 item 流会被交给输出器处理。接下来我们来定义一个输出器管道，默认的输出器在项目配置文件 `piplines.go` 中。(更多输出器的用法见详细介绍部分)

```go

package first

import (
	"github.com/wheat-os/slubby/stream"
	"github.com/wheat-os/wlog"
)

type FirstPipline struct{}

func (t *FirstPipline) OpenSpider() error {
	return nil
}

func (t *FirstPipline) CloseSpider() error {
	return nil
}

func (t *FirstPipline) ProcessItem(item stream.Item) stream.Item {
	wlog.Debug(item)
	return item
}

```

我们执行 `slub crawl`

```
2022-05-27 12:42:24 INFO <Response [200]> Request<url: http://www.baidu.com, method: GET
2022-05-27 12:42:24 DEBUG &{0xc00029e100 百度一下，你就知道}
2022-05-27 12:42:27 INFO the spider shuts down successfully
```

到这里相信你对 `slubby` 已经有了不错的了解了，但是 `slubby` 还有更多有趣的组件等待你的探索，接下来参考，详细文档介绍来获取更多 slubby 的使用方法吧。