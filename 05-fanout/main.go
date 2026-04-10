package main

import (
	"fmt"
	"sync"
)

func fanout(ch <-chan int, splitCount int) []chan int {
	chans := make([]chan int, splitCount)
	for i := range splitCount {
		chans[i] = make(chan int)
	}
	go func() {
		index := 0
		for v := range ch {
			chans[v%len(chans)] <- v
			index++
		}
		for _, c := range chans {
			close(c)
		}
	}()
	return chans
}

func main() {
	ch := make(chan int)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		for i := range 10 {
			ch <- i
		}
		close(ch)
	})

	for i, c := range fanout(ch, 3) {
		wg.Add(1)

		go func(fanOutChan chan int) {
			defer wg.Done()
			for v := range fanOutChan {
				fmt.Printf("Fanout chan %d, value: %d\n", i, v)
			}
		}(c)
	}
	wg.Wait()

}
