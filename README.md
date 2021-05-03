# SF

[![build status](https://img.shields.io/github/workflow/status/0x2E/sf/build)](https://github.com/0x2E/sf/actions/new)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

> **S**ubdomain **F**inder

[English Document](https://github.com/0x2E/sf/blob/main/README_en.md)

SF 是一个高效的子域名收集工具，支持字典爆破等功能，更多模块在不断开发中。

## 安装

三种方式：

1. 在 [release](https://github.com/0x2E/sf/releases) 页面下载编译完成的可执行文件
2. 下载每次 git-push 后自动编译的可执行文件：进入 [Actions](https://github.com/0x2E/sf/actions) 中任意一次 workflow 页面， 下滑页面找到 Artifacts
3. 自行编译 main 分支源码

## 使用

```bash
$ ./sf -h
NAME:
   sf - subdomain finder - https://github.com/0x2E/sf

USAGE:
   sf [global options] [arguments...]

GLOBAL OPTIONS:
   --url value, -u value                        Target url or domain name
   --dict value, -d value                       Load dictionary from a file
   --output value, -o value                     Output results to a file
   --resolver value, -r value                   [fuzz] DNS resolver (default: "8.8.8.8")
   --thread value, -t value                     [fuzz] The number of threads. Each thread will occupy a temporary port of the system until the e
nd of the fuzz (default: 100)
   --queue value, -q value                      [fuzz] The length of the task queue. Too high may fill the system socket buffer and cause packet
 loss (default: 100)
   --wildcardMode value, -w value               [fuzz] Two modes (1 or 2) for processing wildcard records. Mode 1 is only based on the IP blackl
ist. Mode 2 matches the IP blacklist, compares the similarity of web page titles after hits, and degenerates to mode 1 if port 80 cannot be acce
ssed. (default: 1)
   --wildcardBlacklistMaxLen value, --wl value  [fuzz] The maximum length of the IP blacklist for wildcard records (default: 1000)
   --retry value, -R value                      [fuzz] The number of retries (default: 2)
   --help, -h                                   show help (default: false)
```

## TODO

- [开发计划](https://github.com/0x2E/sf/labels/todo)
- [提交新的建议](https://github.com/0x2E/sf/issues/new)
