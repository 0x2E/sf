package fuzz

import (
	"sync"
)

// producer
func producer(ch chan<- string, wg *sync.WaitGroup, dict []string) {
	defer wg.Done()
	for i := range dict {
		ch <- dict[i]
	}
	close(ch)
}
