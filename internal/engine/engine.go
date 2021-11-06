package engine

import (
	"github.com/0x2E/sf/internal/conf"
	"github.com/0x2E/sf/internal/module"
	"github.com/schollz/progressbar/v3"
	"log"
	"sync"
	"time"
)

const (
	RESOLVER = "8.8.8.8"
	THREAD   = 200
	RATE     = 2000
	QUEUELEN = 10000
	RETRY    = 3
	CHECK    = true
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
	valid        []string
	invalid      []string
}

func New(config *conf.Config) *Engine {
	return &Engine{
		conf:  config,
		check: &check{},
		bar:   progressbar.Default(-1, "DNS response receiving"),
		valid: make([]string, 0, QUEUELEN),
		// invalid 在检查existWildcard后创建
		toEnumerator: make(chan *module.Task, QUEUELEN),
		toChecker:    make(chan *module.Task, QUEUELEN),
		toRecorder:   make(chan *module.Task, QUEUELEN),
	}
}

// Run 运行，返回有效和无效结果集，不存在泛解析时无效结果集为nil
func (e *Engine) Run() ([]string, []string) {
	wg := sync.WaitGroup{}
	e.check.existWildcard = e.conf.Check && e.existWildcard()
	if e.check.existWildcard {
		wg.Add(1)
		go e.checker(&wg)
		e.invalid = make([]string, 0, QUEUELEN)
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

	return e.valid, e.invalid
}
