package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if item, ok := c.cache[key]; ok {
		c.ll.MoveToFront(item)
		kv := item.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) Put(key string, value Value) {
	if item, ok := c.cache[key]; ok {
		kv := item.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
		c.ll.MoveToFront(item)
	} else {
		item := c.ll.PushFront(&entry{key, value})
		c.cache[key] = item
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	last_item := c.ll.Back()
	if last_item != nil {
		c.ll.Remove(last_item)
		kv := last_item.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
