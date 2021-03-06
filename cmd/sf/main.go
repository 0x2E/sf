package main

import (
	"github.com/0x2E/sf/internal/engine"
	"github.com/0x2E/sf/internal/option"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
)

const (
	debug   = false
	version = "v0.1.2"
)

func main() {
	if debug {
		//f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		//defer f.Close()
		//pprof.StartCPUProfile(f)
		//defer pprof.StopCPUProfile()

		runtime.GOMAXPROCS(1)              // 限制 CPU 使用数，避免过载
		runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
		runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
		go func() {
			// localhost:10010/debug/pprof
			if err := http.ListenAndServe(":10010", nil); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
	}

	log.SetFlags(log.Ltime)

	o, err := option.ParseOption()
	if err != nil {
		log.Fatal(err)
	}

	e := engine.New(o)
	err = e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
