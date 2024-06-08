package cache

import (
	"fmt"
	"sync"

	"github.com/tom773/cardcache/peer"
)

type Cache struct {
	mu     sync.RWMutex
	Data   map[string][]byte
	PubSub *peer.PubSub
}

func NewCache() *Cache {
	return &Cache{
		Data:   make(map[string][]byte),
		PubSub: peer.NewPubSub(),
	}
}

func (c *Cache) Set(key, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Data == nil {
		c.Data = make(map[string][]byte)
	}
	c.Data[string(key)] = []byte(value)
	c.PubSub.Publish(string(key), string(value))
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

func (c *Cache) Subscribe(p *peer.Peer, key string) chan []byte {
	return c.PubSub.Subscribe(p, key)
}
