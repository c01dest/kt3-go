package main

import (
	"fmt"
	"sync"
)

type ConcurrentMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{m: make(map[string]int)}
}

func (c *ConcurrentMap) Set(key string, val int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = val
}

func (c *ConcurrentMap) Get(key string) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func main() {
	safeMap := NewConcurrentMap()
	var wg sync.WaitGroup
	const total = 100

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			safeMap.Set(fmt.Sprintf("key_%d", id), id*10)
		}(i)
	}

	wg.Wait()
	fmt.Printf("Успешно записано %d элементов в конкурентную map\n", total)
}