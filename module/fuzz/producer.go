package fuzz

import (
	"sync"
)

// producer
func producer(ch chan<- string, wg *sync.WaitGroup, f *FuzzModule) {
	defer wg.Done()
	for i := range f.UnReceived.Data {
		ch <- f.UnReceived.Data[i]
	}
	close(ch)
}
