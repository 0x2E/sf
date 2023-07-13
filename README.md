# SF

[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

中文/[English](https://github.com/0x2E/sf/blob/main/README_en.md)

SF 是一个高效的子域名爆破工具：

- 基于 UDP 的无连接特性并行收发 DNS 请求
- 支持**自定义爆破点**，使用占位符 `%` 设置
- 支持基于 `*` 记录的**泛解析域名检测**
- 支持秒级限流和失败重试
- 支持检测域传送漏洞

## 安装

1. [Release](https://github.com/0x2E/sf/releases)
2. 编译源码

```shell
git clone https://github.com/0x2E/sf.git
cd sf
go build -o sf ./cmd/sf
```
