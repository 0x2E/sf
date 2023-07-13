package main

import (
	"bufio"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/0x2E/sf/internal/conf"
	"github.com/0x2E/sf/internal/engine"
	"github.com/sirupsen/logrus"
)

var Version, Branch, Commit string

func main() {
	var (
		c            = conf.C
		output       string
		disableCheck bool
		slient       bool
		debug        bool
		showHelp     bool
		showVersion  bool
	)
	flag.StringVarP(&c.RawTarget, "domain", "d", "", `Target domain name.
If the placeholder % exists, only replaces the placeholder instead of splicing wordlist as subdomain`)
	flag.StringVarP(&c.Wordlist, "wordlist", "w", "", "Wordlist file")
	flag.StringVarP(&c.Resolver, "resolver", "r", "8.8.8.8", "DNS resolver")
	flag.StringVarP(&output, "output", "o", "", "Output results to a file")
	flag.IntVarP(&c.Concurrent, "concurrent", "t", 800, "Number of concurrent")
	flag.IntVar(&c.Rate, "rate", 30000, `Maximum rate req/s. 
It is recommended to determine if the rate is appropriate by the send/recv statistics in log`)
	flag.IntVar(&c.Retry, "retry", 1, "Number of retries")
	flag.IntVarP(&c.StatisticsInterval, "stats", "s", 2, "Statistics interval(seconds) in log")
	flag.BoolVar(&disableCheck, "disable-check", false, "Disable check the validity of the subdomains")
	flag.BoolVar(&slient, "slient", false, "Only output valid subdomains, and logs that caused abnormal exit, e.g., fatal and panic")
	flag.BoolVar(&debug, "debug", false, "Set the log level to debug, and enable golang pprof with web service")
	flag.BoolVarP(&showVersion, "version", "v", false, "Show version")
	flag.BoolVarP(&showHelp, "help", "h", false, "Show help message")
	flag.CommandLine.SortFlags = false
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		version()
		os.Exit(0)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "20060102 15:04:05",
		FullTimestamp:   true,
	})

	if slient {
		if debug {
			logrus.Fatal("cannot enable 'debug' and 'slient' at the same time")
		}
		logrus.SetLevel(logrus.FatalLevel)
	} else {
		fmt.Print(banner)
		version()
	}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		go pprof()
	}

	c.ValidCheck = !disableCheck

	if err := c.Verify(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("target: [%s]. wordlist: [%s]. resolver: [%s]. concurrent: [%d]. rate: [%d]. retry: [%d]. check valid: [%t]",
		c.RawTarget, c.Wordlist, c.Resolver, c.Concurrent, c.Rate, c.Retry, c.ValidCheck)

	startAt := time.Now()
	res := engine.New().Run()

	logrus.Infof("found %d subdomains. time: %.2f seconds.\n", len(res), time.Since(startAt).Seconds())

	saveResult(output, res)
}

func saveResult(path string, data []string) {
	if strings.TrimSpace(path) == "" || len(data) == 0 {
		return
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o755)
	if err != nil {
		logrus.Error("cannot save results:", err)
		return
	}
	defer f.Close()
	bufWriter := bufio.NewWriter(f)
	for _, v := range data {
		bufWriter.WriteString(v + "\n")
	}
	bufWriter.Flush()
}

func pprof() {
	logrus.Debug("pprof is on 127.0.0.1:10000/debug/pprof")
	if err := http.ListenAndServe("127.0.0.1:10000", nil); err != nil {
		logrus.Error(err)
	}
}

func version() {
	fmt.Printf("version: %s. branch: %s. commit: %s\n", Version, Branch, Commit)
}
