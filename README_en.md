# SF

[![build status](https://img.shields.io/github/workflow/status/0x2E/sf/build)](https://github.com/0x2E/sf/actions/new)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x2E/sf)](https://goreportcard.com/report/github.com/0x2E/sf)
[![go version](https://img.shields.io/github/go-mod/go-version/0x2E/sf)](https://github.com/0x2E/sf/blob/main/go.mod)

[Chinese Document](https://github.com/0x2E/sf/blob/main/README.md)

SF is an efficient subdomain collection tool.

- [x] Brute force
  - Use native socket to send and receive data packets, no need to install additional dependencies.
  - Based on stateless thinking, separate sending and receiving actions make it fast and only occupy a small amount of system resources.
  - Support packet retry function to ensure the stability of the results. 
  - Provides two pan-analysis processing modes
    1. Only based on IP blacklist
    2. Based on IP blacklist and web title similarity
- [x] Domain transfer
- [ ] Third-party database

## Installation

Three ways:

1. Download the compiled executable file on the [release](https://github.com/0x2E/sf/releases) page
2. Download the executable file compiled after each git-push: enter any workflow page in [Actions](https://github.com/0x2E/sf/actions), go down the page to find Artifacts
3. Compile the main branch source code by yourself

## Usage

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

- [Development Plan](https://github.com/0x2E/sf/labels/todo)
- [Submit a new proposal](https://github.com/0x2E/sf/issues/new)
