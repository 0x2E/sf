package engine

import (
	"sync"

	"github.com/0x2E/sf/internal/conf"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

// existWildcard checks if there is a wildcard record
func (e *Engine) existWildcard() bool {
	m := &dns.Msg{}
	m.SetQuestion("*."+conf.C.Target, dns.TypeA)
	resp, err := dns.Exchange(m, conf.C.Resolver)
	if err != nil || resp.Rcode != dns.RcodeSuccess || len(resp.Answer) == 0 {
		return false
	}

	e.wildcardRecord = resp.Answer[0]
	logrus.Debug("found wildcard record: " + e.wildcardRecord.String())
	e.wildcardRecord.Header().Name = "" // for easier comparison
	return true
}

// checker checks if domain is valid:
//
// 1. not wildcard record
func (e *Engine) checker(wg *sync.WaitGroup) {
	defer func() {
		close(e.toRecorder)
		wg.Done()
	}()

	for t := range e.toChecker {
		t.Record.Header().Name = "" // for easier comparison

		if dns.IsDuplicate(t.Record, e.wildcardRecord) {
			continue
		}

		e.toRecorder <- t
	}
}
