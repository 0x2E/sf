# SF

[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

[中文](https://github.com/0x2E/sf/blob/main/README.md)/English

SF is an efficient subdomain brute-forcing tool:

- parallelizes sending and receiving DNS requests based on the connectionless feature of UDP
- supports **custom brute-forcing points** using placeholder `%`
- supports **detecting wildcard record** based on `*` records
- supports second-level rate limiting and retrying on failures
- supports detecting zone-transfer vulnerabilities

## Installation

1. [Release](https://github.com/0x2E/sf/releases)
2. compile

```shell
git clone https://github.com/0x2E/sf.git
cd sf
go build -o sf ./cmd/sf
```
