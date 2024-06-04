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

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Data == nil {
		return nil, fmt.Errorf("Cache Empty")
	}
	value, ok := c.Data[string(key)]
	if !ok {
		return nil, fmt.Errorf("Key not found")
	}
	return value, nil
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
