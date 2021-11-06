package main

import (
	"bufio"
	"flag"
	"github.com/0x2E/sf/internal/conf"
	"github.com/0x2E/sf/internal/engine"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

const (
	DEBUG = false
)

func main() {
	debug()

	var output string
	c := &conf.Config{}
	flag.StringVar(&c.Domain, "u", "", "Target domain name or URL")
	flag.StringVar(&c.Wordlist, "f", "", "Load wordlist from a file")
	flag.StringVar(&c.Resolver, "r", engine.RESOLVER, "DNS resolver")
	flag.StringVar(&output, "o", "", "Output results to a file")
	flag.IntVar(&c.Thread, "t", engine.THREAD, "Number of concurrent")
	flag.IntVar(&c.Rate, "rate", engine.RATE, "Maximum number of DNS requests sent per second")
	flag.IntVar(&c.Retry, "retry", engine.RETRY, "Number of retries")
	flag.BoolVar(&c.Check, "check", engine.CHECK, "Whether to check the validity of the subdomain")
	flag.Parse()

	if err := c.Verify(); err != nil {
		log.Fatal(err)
	}

	if strings.TrimSpace(output) == "" {
		output = c.Domain + "txt" // domain结尾已经有一个点了
	}

	startTime := time.Now()

	app := engine.New(c)
	valid, invalid := app.Run()
	log.Printf("Found %d valid, %d invalid. %.2f seconds in total.\n", len(valid), len(invalid), time.Since(startTime).Seconds())

	saveResult(output, valid)
	saveResult("invalid_"+output, invalid)
}

func saveResult(path string, data []string) {
	if len(data) == 0 {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("cannot save results into file: %s", err)
		return
	}
	defer f.Close()
	bufWriter := bufio.NewWriter(f)
	for _, v := range data {
		bufWriter.WriteString(v + "\n")
	}
	bufWriter.Flush()
	log.Printf("Results are stored in %s\n", path)
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
