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
	for t := range e.toRecorder {
		subdomain = t.Subdomain[:len(t.Subdomain)-1]
		if t.Valid {
			e.valid = append(e.valid, subdomain)
		} else {
			e.invalid = append(e.invalid, subdomain)
		}
	}
	// todo 去重
}
