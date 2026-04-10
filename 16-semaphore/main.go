package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	channel chan struct{}
}

func NewSemaphore(maxConcurrency int) *Semaphore {
	return &Semaphore{
		channel: make(chan struct{}, maxConcurrency),
	}
}

func (s *Semaphore) Acquire() {
	s.channel <- struct{}{}
}

func (s *Semaphore) Release(i int) {
	<-s.channel
}

func (s *Semaphore) doRequest(i int) {
	s.Acquire()
	fmt.Println("Srart doing request ", i)
	time.Sleep(1 * time.Second)
	fmt.Println("End doing request", i)
	s.Release(1)
}

func main() {
	s := NewSemaphore(3)

	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			s.doRequest(i)
			wg.Done()
		}()
	}
	wg.Wait()
}
