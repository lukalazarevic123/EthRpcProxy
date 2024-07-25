package cache

import (
	"container/list"
	"errors"
)

type HolderInfo struct {
	HolderAddress string
	IsHolder      bool
	BlockNumber   int
}

type CacheElement struct {
	element *list.Element
	info    *HolderInfo
}

type LRUCache struct {
	storage  map[string]*CacheElement
	cacheCap int
	ls       list.List
}

func NewLRUCache(cacheCap int) *LRUCache {
	return &LRUCache{
		storage:  map[string]*CacheElement{},
		cacheCap: cacheCap,
		ls:       list.List{},
	}
}

func (c *LRUCache) Get(key string) (*HolderInfo, error) {
	holder, flag := c.storage[key]

	if !flag {
		return nil, errors.New("Element not found")
	}

	c.ls.MoveToFront(holder.element)
	return holder.info, nil
}

func (c *LRUCache) Set(key string, holder *HolderInfo) {
	holderCache, flag := c.storage[key]

	if flag {
		holderCache.info = holder
		c.ls.MoveToFront(holderCache.element)

		return
	}

	newElem := c.ls.PushFront(key)
	c.storage[key] = &CacheElement{
		element: newElem,
		info:    holder,
	}

	if c.ls.Len() > c.cacheCap {
		backElement := c.ls.Back()
		backElementKey := backElement.Value.(string)

		c.ls.Remove(backElement)
		delete(c.storage, backElementKey)
	}
}

func (c *LRUCache) GetAll() []*HolderInfo {
	var allElements []*HolderInfo
	for e := c.ls.Front(); e != nil; e = e.Next() {
		key := e.Value.(string)
		if element, found := c.storage[key]; found {
			allElements = append(allElements, element.info)
		}
	}
	return allElements
}

func (c *LRUCache) GetKeys() []string {
	keys := make([]string, 0, len(c.storage))
	for key := range c.storage {
		keys = append(keys, key)
	}
	return keys
}
