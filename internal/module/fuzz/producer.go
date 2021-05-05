package fuzz

import (
	"sync"
)

// producer
func producer(ch chan<- string, wg *sync.WaitGroup, f *FuzzModule) {
	defer wg.Done()
	for i := range f.unReceived.data {
		ch <- f.unReceived.data[i]
	}
	close(ch)
}
