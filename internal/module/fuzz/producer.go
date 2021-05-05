package fuzz

import (
	"sync"
)

// producer
func producer(ch chan<- string, wg *sync.WaitGroup, unReceived []string) {
	defer wg.Done()
	for i := range unReceived {
		ch <- unReceived[i]
	}
	close(ch)
}
