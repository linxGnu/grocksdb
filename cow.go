package grocksdb

import (
	"sync"
	"sync/atomic"
)

// COWList implements a copy-on-write list. It is intended to be used by go
// callback registry for CGO, which is read-heavy with occasional writes.
// Reads do not block; Writes do not block reads (or vice versa), but only
// one write can occur at once;
type COWList struct {
	v  *atomic.Value
	mu *sync.Mutex
}

// NewCOWList creates a new COWList.
func NewCOWList() *COWList {
	var list []interface{}
	v := &atomic.Value{}
	v.Store(list)
	return &COWList{v: v, mu: new(sync.Mutex)}
}

// Append appends an item to the COWList and returns the index for that item.
func (c *COWList) Append(i interface{}) (index int) {
	c.mu.Lock()
	list := c.v.Load().([]interface{})
	newLen := len(list) + 1
	newList := make([]interface{}, newLen)
	copy(newList, list)
	newList[newLen-1] = i
	c.v.Store(newList)
	c.mu.Unlock()
	index = newLen - 1
	return
}

// Get gets the item at index.
func (c *COWList) Get(index int) interface{} {
	list := c.v.Load().([]interface{})
	return list[index]
}
