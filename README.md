# SF

[![build status](https://img.shields.io/github/actions/workflow/status/0x2E/sf/build.yml?branch=main)](https://github.com/0x2E/sf/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

SF 是一个高效的子域名收集工具。目前已有字典爆破、域传送模块。

<details>
    <summary>演示</summary>
    <a href="https://asciinema.org/a/447397" target="_blank"><img src="https://asciinema.org/a/447397.svg" /></a>
</details>

子域名由各个模块收集后送入任务队列，按需进行域名解析（Enumerator）、有效性检测（Checker）、记录（Recorder）。

Enumerator 基于 UDP 的无状态特性，将发送和接收分离，效率更高，支持限流和重试机制。

Checker 根据 DNS 记录等特征筛选有效子域名，目前用于泛解析检测（[关于泛解析检测的一些问题](https://github.com/0x2E/sf/issues/12)）。


## 安装

- 使用编译好的可执行文件
  - 稳定：[release](https://github.com/0x2E/sf/releases)
  - 主分支最新：进入 [Actions](https://github.com/0x2E/sf/actions) 中任意一次 workflow，下滑页面找到 Artifacts
- 编译源码

## 使用方法

- `-u` 目标域名（必需）
- `-f` 字典路径，为空时不启动爆破模块
- `-r` DNS 服务器，默认 `8.8.8.8`
- `-o` 结果输出路径，默认 `{domain}.txt`
- `-check` 是否开启有效性检查，默认 true
- `-t` 并发数，默认 200
- `-rate` 每秒最大请求量，默认 2000
- `-retry` 重试次数，默认 3
- `-h` 输出完整参数列表

## TODO

- [开发计划](https://github.com/0x2E/sf/labels/todo)
- [提出建议](https://github.com/0x2E/sf/issues/new)
