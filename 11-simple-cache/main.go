package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func (c *Cache) longCalculation(ctx context.Context, n int) (int, error) {
	ch := make(chan int, 1)

	go func() {
		secondsToSleep := rand.Float64() * float64(n)
		time.Sleep(time.Duration(secondsToSleep) + 2*time.Second)
		ch <- n + 1
		defer close(ch)
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case v := <-ch:
		return v, nil
	}
}

type Cache struct {
	cache   map[int]int
	mu      sync.RWMutex
	hits    atomic.Int64
	misses  atomic.Int64
	timeout time.Duration
}

func NewCache(defaultTimeout time.Duration) *Cache {
	return &Cache{
		cache:   make(map[int]int),
		timeout: defaultTimeout,
	}
}
func (c *Cache) HitRate() float64 {
	return 1.0 - (float64(c.misses.Load()) / float64(c.hits.Load()))
}

func (c *Cache) CachedLongCalculation(n int) (int, error) {
	c.hits.Add(1)
	c.mu.RLock()
	found, ok := c.cache[n]
	c.mu.RUnlock()

	if !ok {
		c.misses.Add(1)
		c.mu.Lock()
		defer c.mu.Unlock()

		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
		value, err := c.longCalculation(ctx, n)
		if err != nil {
			return 0, err
		}

		c.cache[n] = value

		return value, nil
	}
	return found, nil
}

func main() {
	nums := []int{5, 10, 22, 234, 234, 5432, 6, 74, 5, 5, 10, 22, 234}
	cache := NewCache(1 * time.Second)
	var wg sync.WaitGroup

	for _, n := range nums {
		wg.Go(func() {
			val, err := cache.CachedLongCalculation(n)
			if err != nil {
				fmt.Println(err)
			} else {

				fmt.Printf("LongCaclcuation(%d) = %d\n", n, val)
			}
		})
	}
	wg.Wait()
	fmt.Println(cache.HitRate())

	fmt.Printf("Goroutins: %d\n", runtime.NumGoroutine())
}
