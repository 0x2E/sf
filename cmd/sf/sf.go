package main

import (
	"flag"
	"github.com/0x2E/sf/internal/conf"
	"github.com/0x2E/sf/internal/engine"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const (
	DEBUG = false
)

func main() {
	debug()

	config := parseArgs()
	err := engine.New(config).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func parseArgs() *conf.Config {
	c := &conf.Config{}
	flag.StringVar(&c.Domain, "u", "", "Target domain name or URL")
	flag.StringVar(&c.Wordlist, "f", "", "Load wordlist from a file")
	flag.StringVar(&c.Resolver, "r", engine.RESOLVER, "DNS resolver")
	flag.StringVar(&c.Output, "o", "", "Output results to a file")
	flag.IntVar(&c.Thread, "t", engine.THREAD, "Number of concurrent")
	flag.IntVar(&c.Rate, "rate", engine.RATE, "Maximum number of DNS requests sent per second")
	flag.IntVar(&c.Retry, "retry", engine.RETRY, "Number of retries")
	flag.BoolVar(&c.Check, "check", engine.CHECK, "Whether to check the validity of the subdomain")
	flag.Parse()

	if err := c.Verify(); err != nil {
		log.Fatal(err)
	}
	return c
}

func debug() {
	if DEBUG {
		//f, _ := os.OpenFile("sf-dev/cpu.pprof", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		//defer f.Close()
		//pprof.StartCPUProfile(f)
		//defer pprof.StopCPUProfile()

		go func() {
			// localhost:10010/debug/pprof
			if err := http.ListenAndServe(":10010", nil); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
	}
}
