package cache

import (
	"sync"
	"time"
)

type Entry[T any] struct {
	Value     T
	ExpiresAt time.Time
}

type Cache[T any] struct {
	mu      sync.Mutex
	entries map[string]Entry[T]
}

func New[T any]() *Cache[T] {
	return &Cache[T]{
		entries: make(map[string]Entry[T]),
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[key]

	if !exists || time.Now().After(entry.ExpiresAt) {
		var zero T
		return zero, false
	}

	return entry.Value, true
}

func (c *Cache[T]) Set(
	key string,
	value T,
	duration time.Duration,
) {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = Entry[T]{
		Value:     value,
		ExpiresAt: time.Now().Add(duration),
	}
}
