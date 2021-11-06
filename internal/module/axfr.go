package module

import (
	"github.com/0x2E/sf/internal/conf"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"sync"
)

// Axfr 域传送模块
//
// 使用zonetransfer.me测试（https://digi.ninja/projects/zonetransferme.php）
type Axfr struct {
	base
}

func newAxfr(conf *conf.Config, toRecorder chan<- *Task) *Axfr {
	return &Axfr{
		base: base{
			name:   "zone-transfer",
			conf:   conf,
			toNext: toRecorder,
		},
	}
}

func (a *Axfr) Run() error {
	m := new(dns.Msg)
	m.SetQuestion(a.conf.Domain, dns.TypeNS)
	r, err := dns.Exchange(m, a.conf.Resolver)
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
		go a.doTransfer(&wg, n.Ns)
	}
	wg.Wait()
	return nil
}

func (a *Axfr) doTransfer(wg *sync.WaitGroup, ns string) {
	defer wg.Done()

	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(a.conf.Domain)
	recvChan, err := t.In(m, ns+":53") // 默认2秒超时
	if err != nil {
		return
	}
	for v := range recvChan {
		if v.Error != nil { //todo 完善对错误类型的判断
			break
		}
		for _, rr := range v.RR {
			t := rr.Header().Rrtype
			if t == dns.TypeA || t == dns.TypeAAAA || t == dns.TypeCNAME || t == dns.TypeMX { //todo 补充类型
				NewTask(a.toNext, rr.Header().Name)
			}
		}
	}
}
