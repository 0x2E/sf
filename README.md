# SF

[![build status](https://img.shields.io/github/workflow/status/0x2E/sf/build)](https://github.com/0x2E/sf/actions/new)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

> **S**ubdomain **F**inder

[English Document](https://github.com/0x2E/sf/blob/main/README_en.md)

SF 是一个高效的子域名收集工具，支持字典爆破等功能，更多模块在不断开发中。

## 安装

你可以在 [release](https://github.com/0x2E/sf/releases) 下载已经编译好的二进制文件，或者使用 main 分支代码自行编译。

## 使用

```bash
./sf -u example.com
```

|标志|功能|默认值|
|:-:|:-:|:-:|
|u|「url」目标域名||
|d|「dict」爆破字典|./dict.txt|
|o|「output」保存结果的文件|./{{ 域名 }}.{{ 时间 }}.txt|
|r|「resolver」DNS 解析服务器|8.8.8.8|
|t|「thread」并发数|100|
|q|「queue」UDP 发送-接收队列长度|100|
|w|「wildcard」泛解析处理模式：简易模式 1，严格模式 2|1|
|R|「retry」重试次数|2|

## TODO

- [开发计划](https://github.com/0x2E/sf/labels/todo)
- [提交新的建议](https://github.com/0x2E/sf/issues/new)
