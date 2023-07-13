package module

import (
	"context"
	"sync"

	"github.com/0x2E/sf/internal/conf"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// RunAxfr is zone transfer module
//
// test: https://digi.ninja/projects/zonetransferme.php
func RunAxfr(ctx context.Context, toNext chan<- *Task) error {
	m := new(dns.Msg)
	m.SetQuestion(conf.C.Target, dns.TypeNS)
	r, err := dns.Exchange(m, conf.C.Resolver)
	if err != nil {
		return errors.Wrap(err, " get NS record")
	}
	if r.Rcode != dns.RcodeSuccess || len(r.Answer) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	var n *dns.NS
	var ok bool
	for _, v := range r.Answer {
		if n, ok = v.(*dns.NS); !ok {
			continue
		}
		wg.Add(1)
		go transferOneNS(&wg, n.Ns, toNext)
	}
	wg.Wait()
	return nil
}

func transferOneNS(wg *sync.WaitGroup, ns string, toNext chan<- *Task) {
	defer wg.Done()

	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(conf.C.Target)
	recvChan, err := t.In(m, ns+":53") // default timeout 2s
	if err != nil {
		return
	}
	for v := range recvChan {
		if v.Error != nil { // TODO: more error type
			break
		}
		for _, rr := range v.RR {
			t := rr.Header().Rrtype
			if t == dns.TypeA || t == dns.TypeAAAA || t == dns.TypeCNAME || t == dns.TypeMX { // TODO: more type
				putTask(toNext, rr.Header().Name)
			}
		}
	}
}
