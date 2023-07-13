package engine

import (
	"sync"

	"github.com/0x2E/sf/internal/conf"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

// existWildcard checks if there is a wildcard record
func (e *Engine) existWildcard() bool {
	m := new(dns.Msg)
	m.SetQuestion(conf.C.Target, dns.TypeNS)
	r, err := dns.Exchange(m, conf.C.Resolver)
	if err != nil || r.Rcode != dns.RcodeSuccess || len(r.Answer) == 0 {
		return false
	}
	for _, v := range r.Answer {
		n, ok := v.(*dns.NS)
		if !ok {
			continue
		}
		m := &dns.Msg{}
		m.SetQuestion("*."+conf.C.Target, dns.TypeA)
		resp, err := dns.Exchange(m, n.Ns+":53")
		if err != nil || resp.Rcode != dns.RcodeSuccess || len(resp.Answer) == 0 {
			continue
		}
		e.wildcardRecord = resp.Answer
		break
	}
	if len(e.wildcardRecord) == 0 {
		return false
	}
	e.wildcardRecord[0].Header().Name = "" // for easier comparison

	return true
}

// checker checks if domain is valid
//
// more: https://github.com/0x2E/sf/issues/12
func (e *Engine) checker(wg *sync.WaitGroup) {
	defer func() {
		close(e.toRecorder)
		wg.Done()
	}()

	logger := logrus.WithField("step", "checker")

	for t := range e.toChecker {
		if len(t.Answer) == len(e.wildcardRecord) {
			matchAll := true
			t.Answer[0].Header().Name = "" // for easier comparison
			for i, v := range e.wildcardRecord {
				if !dns.IsDuplicate(v, t.Answer[i]) {
					matchAll = false
					break
				}
			}
			if matchAll {
				t.Valid = false
				logger.Debug("invalid: " + t.DomainName)
			}
		}
		e.toRecorder <- t
	}
}
