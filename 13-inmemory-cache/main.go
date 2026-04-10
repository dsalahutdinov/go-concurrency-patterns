package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Cache interface {
	Set(k string, v string)
	Get(k string) (string, bool)
}

type ShardedCache struct {
	shards []Cache
}

func (sc *ShardedCache) Set(k string, v string) {
	sc.shards[sc.getShardIndex(k)].Set(k, v)

}

func (sc *ShardedCache) Get(k string) (string, bool) {
	return sc.shards[sc.getShardIndex(k)].Get(k)
}

func NewShardedCache(size int) Cache {
	shards := make([]Cache, size)
	for i := range size {
		shards[i] = NewSingleCache()
	}
	return &ShardedCache{
		shards: shards,
	}
}

func (sc *ShardedCache) getShardIndex(k string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(k))
	return h.Sum64() % uint64(len(sc.shards))
}

type SingleCache struct {
	data map[string]string
	mu   sync.RWMutex
}

func (sc *SingleCache) Set(k string, v string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.data[k] = v
}

func (sc *SingleCache) Get(k string) (string, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	value, ok := sc.data[k]
	return value, ok
}

func NewSingleCache() Cache {
	return &SingleCache{
		data: make(map[string]string),
	}
}

func main() {
	var cache Cache = NewShardedCache(3)

	testData := []string{"asdf", "sdfg", "dfgh", "dfgh"}
	var wg sync.WaitGroup

	for i := range 1000 {
		wg.Go(func() {
			v := testData[i%len(testData)]
			cache.Set(v, v)

		})
	}
	wg.Wait()
	fmt.Println(cache.Get("asdf"))

}
