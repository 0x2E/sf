# SF

[![build status](https://img.shields.io/github/workflow/status/0x2E/sf/build)](https://github.com/0x2E/sf/actions/new)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

> **S**ubdomain **F**inder

[Chinese Document](https://github.com/0x2E/sf/blob/main/README.md)

SF is an efficient subdomain collection tool that supports dictionary blasting and other functions. More modules are under development.

## Installation

You can download the compiled binary file at [release](https://github.com/0x2E/sf/releases), or use the main branch code to compile it yourself.

## Usage

```bash
./sf -u example.com
```

|flags|function|default|
|:-:|:-:|:-:|
|u|「url」target domain name||
|d|「dict」dictionary|./dict.txt|
|o|「output」file to save results|./{{ domain-name }}.{{ start-time }}.txt|
|r|「resolver」DNS resolver|8.8.8.8|
|t|「thread」number of thread|100|
|q|「queue」UDP send-receive queue length|100|
|w|「wildcard」wildcard processing mode: simple mode 1, strict mode 2|1|
|R|「retry」number of retries|2|

## TODO

- [Development Plan](https://github.com/0x2E/sf/labels/todo)
- [Submit a new proposal](https://github.com/0x2E/sf/issues/new)
