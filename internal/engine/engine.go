package engine

import (
	"context"
	"sync"
	"time"

	"github.com/0x2E/sf/internal/conf"
	"github.com/0x2E/sf/internal/module"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

const (
	QueueMaxLen = 1000
)

var (
	NetTimeout = 3 * time.Second
)

type Engine struct {
	needCheck bool
	// wildcardRecord is the wildcard record (`*.example.com`)
	wildcardRecord []dns.RR
	toResolver     chan *module.Task
	toChecker      chan *module.Task
	toRecorder     chan *module.Task
	validResults   []string
	invalidResults []string
}

func New(config *conf.Config) *Engine {
	return &Engine{
		toResolver: make(chan *module.Task, QueueMaxLen),
		toChecker:  make(chan *module.Task, QueueMaxLen),
		toRecorder: make(chan *module.Task, QueueMaxLen),
	}
}

func (e *Engine) Run() ([]string, []string) {
	wg := sync.WaitGroup{}
	e.needCheck = conf.C.ValidCheck && e.existWildcard()
	if e.needCheck {
		wg.Add(1)
		go e.checker(&wg)
		logrus.Debugf("wirldcard record: %#v", e.wildcardRecord)
	} else {
		logrus.Debug("turn off checker")
		close(e.toChecker)
	}
	wg.Add(2)
	go e.resolver(&wg)
	go e.recorder(&wg)

	wgModules := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wgModules.Add(1)
	go func() {
		startAt := time.Now()
		logger := logrus.WithField("module", "zone-transfer")
		defer wgModules.Done()

		if err := module.RunAxfr(ctx, e.toRecorder); err != nil {
			logger.Error(err)
		}

		logger.Debug("done, time: " + time.Since(startAt).String())
	}()
	wgModules.Add(1)
	go func() {
		startAt := time.Now()
		logger := logrus.WithField("module", "wordlist")
		defer wgModules.Done()

		if err := module.RunWordlist(ctx, e.toResolver); err != nil {
			logger.Error(err)
		}

		logger.Debug("done, time: " + time.Since(startAt).String())
	}()

	wgModules.Wait()
	logrus.Debug("all modules done")

	close(e.toResolver)
	wg.Wait()

	return e.validResults, e.invalidResults
}
