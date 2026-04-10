package main

import (
	"context"
	"log"
	"time"
)

var counter int64

func SimulateRequest(ctx context.Context) (int64, error) {
	ch := make(chan int64)
	defer close(ch)

	go func() {
		time.Sleep(time.Duration(3) * time.Second)

		counter++
		ch <- counter
	}()

	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	case val := <-ch:
		return val, nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	val, err := SimulateRequest(ctx)
	log.Printf("Значение счетчика: %d, err: %s\n", val, err)
}
