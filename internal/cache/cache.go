package cache

import (
	"sync"

	"order-service-wb/internal/models"
)

type Cache struct {
	mu    *sync.RWMutex
	store map[string]models.Order
	order []string
	size  int
}

func NewCache(size int) *Cache {
	return &Cache{
		store: make(map[string]models.Order),
		order: make([]string, 0, size),
		size:  size,
	}
}

func (c *Cache) Set(id string, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.store[id]; ok {
		c.store[id] = order
		return
	}

	if len(c.order) >= c.size {
		old := c.order[0]
		delete(c.store, old)
		c.order = c.order[1:]
	}

	c.store[id] = order
	c.order = append(c.order, id)
}

func (c *Cache) Get(id string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[id]
	return val, ok
}
