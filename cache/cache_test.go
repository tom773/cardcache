package cache

import (
	"fmt"
	"testing"
)

var c = NewCache()

func TestSetCache(t *testing.T) {
	key := []byte("foo")
	value := []byte("bar\n")

	err := c.Set(key, value)
	if err != nil {
		t.Fatalf("Error setting key: %v", err)
	}
	fmt.Printf("Set key %s with value %s\n", string(key), string(value))
}

func TestGetCache(t *testing.T) {
	key := []byte("foo")
	value, err := c.Get(key)
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}
	fmt.Printf("Value for key %s: %s\n", string(key), string(value))
}

func TestDelCache(t *testing.T) {
	key := []byte("foo")
	err := c.Del(key)
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
	fmt.Printf("Deleted key %s\n", string(key))
}

func TestGetCacheEmpty(t *testing.T) {
	key := []byte("foo")
	_, err := c.Get(key)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	fmt.Printf("Cache is empty\n")
}
