package module

import (
	"github.com/miekg/dns"
)

type Task struct {
	DomainName  string
	Record      dns.RR
	LastQueryAt int64
	Received    bool
}

func putTask(toNext chan<- *Task, dn string) {
	toNext <- &Task{
		DomainName: dns.Fqdn(dn),
	}
}
