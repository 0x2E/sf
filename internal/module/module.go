package module

import (
	"github.com/miekg/dns"
)

type Task struct {
	DomainName  string
	Answer      []dns.RR
	LastQueryAt int64
	Received    bool
	Valid       bool
}

func putTask(toNext chan<- *Task, dn string) {
	toNext <- &Task{
		DomainName: dns.Fqdn(dn),
		Valid:      true,
	}
}
