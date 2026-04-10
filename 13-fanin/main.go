package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	channels := make([]chan int64, 100000)

	for i := range channels {
		channels[i] = make(chan int64)
	}

	for i := range channels {
		go func(i int) {
			channels[i] <- int64(i)
			close(channels[i])
		}(i)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	for v := range merge2(ctx, channels...) {
		fmt.Println(v)
	}
}

// func merge(channels ...chan int64) chan int64 {
// 	ch := make(chan int64)

// 	var wg sync.WaitGroup

// 	for _, c := range channels {
// 		wg.Add(1)
// 		go func(source chan int64) {
// 			for v := range source {
// 				ch <- v
// 			}
// 			wg.Done()
// 		}(c)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(ch)
// 	}()

// 	return ch
// }

func merge2(ctx context.Context, channels ...chan int64) chan int64 {
	ch := make(chan int64)

	var wg sync.WaitGroup

	for _, c := range channels {
		wg.Add(1)
		go func(source chan int64) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					fmt.Println("context done")
					return
				case val, ok := <-source:
					if !ok {
						return
					}
					ch <- val
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
