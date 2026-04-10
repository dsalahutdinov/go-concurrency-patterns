package main

import (
	"fmt"
	"sync"
)

type Pool struct {
	jobs chan int
}

func NewPool() *Pool {
	return &Pool{}
}

func (p *Pool) work(jobs <-chan int) <-chan int {
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for x := range jobs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				results <- x * x
			}()
		}
		defer wg.Done()
	}()
	go func() {
		defer close(results)

		wg.Wait()
	}()
	return results
}

func main() {

	pool := NewPool()

	jobs := make(chan int)

	results := pool.work(jobs)

	go func() {
		for i := range 10 {
			jobs <- i
		}
		close(jobs)
	}()

	for r := range results {
		fmt.Println(r)
	}
}
