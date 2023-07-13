package engine

import (
	"fmt"
	"sync"
)

func (e *Engine) recorder(wg *sync.WaitGroup) {
	defer wg.Done()

	validSet, invalidSet := make(map[string]struct{}), make(map[string]struct{})
	for t := range e.toRecorder {
		subdomain := t.DomainName[:len(t.DomainName)-1]
		if t.Valid {
			fmt.Println(subdomain)
			validSet[subdomain] = struct{}{}
		} else {
			invalidSet[subdomain] = struct{}{}
		}
	}

	e.validResults = make([]string, 0, len(validSet))
	for d := range validSet {
		e.validResults = append(e.validResults, d)
	}
	e.invalidResults = make([]string, 0, len(invalidSet))
	for d := range invalidSet {
		e.invalidResults = append(e.invalidResults, d)
	}
}
