package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type CacheShard interface {
	Set(k string, v string)
	Get(k string) (bool, string)
}

type ShardedCacheImpl struct {
	shardsCount int
	shards      []CacheShard
}

func (c *ShardedCacheImpl) Set(k string, v string) {
	s, err := c.getShard(k)
	if err != nil {

		panic("No shard found")
	}

	s.Set(k, v)
}

func (c *ShardedCacheImpl) Get(k string) (bool, string) {
	s, err := c.getShard(k)
	if err == nil {
		b, r := s.Get(k)
		return b, r
	} else {
		return false, ""
	}
}

func (c *ShardedCacheImpl) getShard(k string) (CacheShard, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(k))
	if err != nil {
		return nil, err
	}

	sum := h.Sum64()
	index := sum % uint64(c.shardsCount)
	return c.shards[index], nil

}

func NewShardedCacheImpl(shardsCount int) *ShardedCacheImpl {
	shards := make([]CacheShard, 0, 5)
	for range shardsCount {
		shards = append(shards, NewCacheShardImpl())
	}

	return &ShardedCacheImpl{
		shards:      shards,
		shardsCount: shardsCount,
	}
}

type CacheShardImpl struct {
	m    sync.RWMutex
	hash map[string]string
}

func NewCacheShardImpl() (c *CacheShardImpl) {
	return &CacheShardImpl{
		hash: make(map[string]string),
	}
}

func (c *CacheShardImpl) Set(k string, v string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.hash[k] = v
}

func (c *CacheShardImpl) Get(k string) (bool, string) {
	c.m.RLock()
	defer c.m.RUnlock()

	value, ok := c.hash[k]
	return ok, value
}

func main() {
	c := NewShardedCacheImpl(5)

	c.Set("foo", "bar")
	ok, s := c.Get("foo")
	fmt.Println(ok, s)
}
