# SF

[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

SF 是一个高效的子域名爆破工具：

- 基于 UDP 的无连接特性并行收发 DNS 请求
- 支持检测域传送漏洞
- 支持秒级限流和失败重试
- 支持基于 `*.` 记录的泛解析域名检测

## 安装

1. [Release](https://github.com/0x2E/sf/releases)
2. 编译源码

```shell
git clone https://github.com/0x2E/sf.git
cd sf
go build -o sf ./cmd/sf
```
