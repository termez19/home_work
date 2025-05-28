package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type pair struct {
	key   Key
	value any
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		item.Value = pair{key, value}
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() == c.capacity {
		back := c.queue.Back()
		if back != nil {
			p := back.Value.(pair)
			println("[Set] Evicting key:", string(p.key))
			delete(c.items, p.key)
			c.queue.Remove(back)
		}
	}

	println("[Set] Inserting new key:", string(key))
	item := c.queue.PushFront(pair{key, value})
	c.items[key] = item
	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(pair).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
