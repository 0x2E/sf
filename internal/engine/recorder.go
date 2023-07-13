package engine

import (
	"fmt"
	"sync"
)

func (e *Engine) recorder(wg *sync.WaitGroup) {
	defer wg.Done()

	res := make(map[string]struct{})
	for t := range e.toRecorder {
		subdomain := t.DomainName[:len(t.DomainName)-1]
		fmt.Println(subdomain)
		res[subdomain] = struct{}{}
	}

	e.results = make([]string, 0, len(res))
	for d := range res {
		e.results = append(e.results, d)
	}
}
