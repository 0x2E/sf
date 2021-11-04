package engine

import (
	"bufio"
	"github.com/0x2E/sf/internal/conf"
	"github.com/0x2E/sf/internal/module"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"sync"
	"time"
)

const (
	RESOLVER = "8.8.8.8"
	THREAD   = 200
	RATE     = 2000
	QUEUELEN = 10000
	RETRY    = 3
	CHECK    = false
)

var (
	NETTIMEOUT = 3 * time.Second
)

type Engine struct {
	conf         *conf.Config
	check        *check
	bar          *progressbar.ProgressBar
	toEnumerator chan *module.Task
	toChecker    chan *module.Task
	toRecorder   chan *module.Task
	result       map[string]struct{}
}

func New(config *conf.Config) *Engine {
	return &Engine{
		conf:         config,
		check:        &check{},
		bar:          progressbar.Default(-1, "DNS response receiving"),
		result:       make(map[string]struct{}),
		toEnumerator: make(chan *module.Task, QUEUELEN),
		toChecker:    make(chan *module.Task, QUEUELEN),
		toRecorder:   make(chan *module.Task, QUEUELEN),
	}
}

func (e *Engine) Run() error {
	startTime := time.Now()
	wg := sync.WaitGroup{}
	e.check.existWildcard = e.conf.Check && e.existWildcard()
	if e.check.existWildcard {
		wg.Add(1)
		go e.checker(&wg)
	} else {
		close(e.toChecker)
	}
	wg.Add(2)
	go e.enumerator(&wg)
	go e.recorder(&wg)

	modules := module.Load(e.conf, e.toEnumerator, e.toRecorder)
	wgModules := sync.WaitGroup{}
	for i := range modules {
		wgModules.Add(1)
		go func(m module.Module) {
			defer wgModules.Done()
			err := m.Run()
			if err != nil {
				log.Printf("[%s] error: %s\n", m.GetName(), err)
				return
			}
		}(modules[i])
	}
	wgModules.Wait()
	close(e.toEnumerator)
	wg.Wait()
	e.bar.Finish()

	log.Printf("Found %d valid subdomains. %s seconds in total.\n", len(e.result), time.Since(startTime))

	if len(e.result) == 0 {
		return nil
	}
	f, err := os.OpenFile(e.conf.Output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "create output file")
	}
	defer f.Close()

	bufWriter := bufio.NewWriter(f)
	for k := range e.result {
		_, _ = bufWriter.WriteString(k + "\n")
	}
	_ = bufWriter.Flush()
	log.Printf("results are stored in %s", e.conf.Output)
	return nil
}
