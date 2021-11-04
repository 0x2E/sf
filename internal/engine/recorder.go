package engine

import (
	"sync"
)

func (e *Engine) recorder(wg *sync.WaitGroup) {
	defer func() {
		//log.Println("Recorder done.")
		wg.Done()
	}()

	var subdomain string
	var ok bool
	for t := range e.toRecorder {
		subdomain = t.Subdomain[:len(t.Subdomain)-1]
		if _, ok = e.result[subdomain]; ok {
			continue
		}
		e.result[subdomain] = struct{}{}
	}
}
