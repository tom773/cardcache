package cache

import (
	"fmt"
	"sync"
)

type Cache struct {
	mu   sync.RWMutex
	Data map[string][]byte
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Set(key, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Data == nil {
		c.Data = make(map[string][]byte)
	}
	c.Data[string(key)] = []byte(value)
	return nil
}

func (c *Cache) Get(key []byte) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Data == nil {
		return []byte("Cache Empty"), true
	}
	value, ok := c.Data[string(key)]
	if !ok {
		return []byte("Key not found"), true
	}
	return value, false
}

func (c *Cache) Del(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Data == nil {
		return fmt.Errorf("Key not found")
	}
	_, ok := c.Data[string(key)]
	if !ok {
		return fmt.Errorf("Key not found")
	}
	delete(c.Data, string(key))
	return nil
}
