# SF

[![build status](https://img.shields.io/github/workflow/status/0x2E/sf/build)](https://github.com/0x2E/sf/actions/new)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

> **S**ubdomain **F**inder

[Chinese Document](https://github.com/0x2E/sf/blob/main/README.md)

SF is an efficient subdomain collection tool that supports dictionary blasting and other functions. More modules are under development.

## Installation

Three ways:

1. Download the compiled executable file on the [release](https://github.com/0x2E/sf/releases) page
2. Download the executable file compiled after each git-push: enter any workflow page in [Actions](https://github.com/0x2E/sf/actions), go down the page to find Artifacts
3. Compile the main branch source code by yourself

## Usage

|flags|function|default|
|:-:|:-:|:-:|
|u|「url」target domain name||
|d|「dict」dictionary path|[built-in dictionary](https://github.com/0x2e/sf/blob/main/module/fuzz/dict.txt)|
|o|「output」output path|[domain name]-[timestamp].txt|
|r|「resolver」DNS resolver|8.8.8.8|
|t|「thread」number of thread|100|
|q|「queue」UDP send-receive queue length|100|
|w|「wildcard」wildcard processing mode: simple mode 1, strict mode 2|1|
|R|「retry」number of retries|2|

## TODO

- [Development Plan](https://github.com/0x2E/sf/labels/todo)
- [Submit a new proposal](https://github.com/0x2E/sf/issues/new)
