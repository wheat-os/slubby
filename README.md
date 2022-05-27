# slubby

[English documentation](./README.en.md)

### slubby 简介

slubby 是一个 基于 go 语言的、组件化的、高扩展性的、快速的、爬虫开发组件库，可以配合 slub 手脚架工具实现 go 语言爬虫的快速快发。

### 功能特性

- 组件化，可以自由替换支持组件来获取不同的表现。
  
- 扩展性良好，支持多个过程中间件，高效扩展爬虫功能。
  
- 搭配 slub 手脚架快速开发爬虫。
  
- 更多常用爬虫组件支持。
  

### 要求

- go 1.15
  
- go mod
  

### 安装方法（slub）

```shell
# 这个方法会安装 slub 爬虫手脚架，我们将使用 slub 来创建 slubby 项目
go install github/wheat-os/slub@latest
```

### 快速开始
[查看快速开始文档](./docs/use/first.md)