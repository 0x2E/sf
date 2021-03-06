# SF

[![build status](https://img.shields.io/github/workflow/status/0x2E/sf/build)](https://github.com/0x2E/sf/actions/new)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

[English Document](https://github.com/0x2E/sf/blob/main/README_en.md)

SF 是一个高效的子域名收集工具。

- [x] 字典爆破
  - 使用原生 socket 收发数据包，无需安装额外依赖
  - 基于无状态爆破思想，发送/接收分离，速度快且仅占用少量系统资源
  - 支持丢包重试机制，保障结果稳定性
  - 提供两种泛解析处理模式
    1. 仅基于 IP 黑名单
    2. 基于 IP 黑名单和网页标题相似度
- [x] 域传送
- [ ] 第三方数据库

## 安装

三种方式：

1. 在 [release](https://github.com/0x2E/sf/releases) 页面下载编译完成的可执行文件
2. 下载每次 git-push 后自动编译的可执行文件：进入 [Actions](https://github.com/0x2E/sf/actions) 中任意一次 workflow 页面， 下滑页面找到 Artifacts
3. 自行编译 main 分支源码

## 使用

```bash
$ ./sf -h
Usage of ./sf:
  -R int
        [fuzz] The number of retries (default 2)
  -d string
        Load dictionary from a file
  -o string
        Output results to a file
  -q int
        [fuzz] The length of the task queue. Too high may fill the system socket buffer and cause packet loss (default 100)
  -r string
        [fuzz] DNS resolver (default "8.8.8.8")
  -t int
        [fuzz] The number of threads. Each thread will occupy a temporary port of the system until the end of the fuzz (default 100)
  -u string
        Target url or domain name
  -w int
        [fuzz] Two modes (1 or 2) for processing wildcard records. Mode 1 is only based on the IP blacklist. Mode 2 matches the IP blacklist, compares the similarity of web page titles after hits, and degenerates to mode 1 if port 80 cannot be accessed (default 1)
  -wl int
        [fuzz] The maximum length of the IP blacklist for wildcard records (default 1000)
```

## TODO

- [开发计划](https://github.com/0x2E/sf/labels/todo)
- [提出建议](https://github.com/0x2E/sf/issues/new)
