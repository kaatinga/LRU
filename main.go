package LRU

import (
	"errors"
	"sync"
)

// The LRU Cache Item Model
type item struct {
	count    byte
	next     *item
	previous *item
	index    string
	data     interface{}
}

// The LRU Cache Order SubModel
type order struct {
	head *item
	tail *item
}

// The LRU Cache Model
type Cache struct {
	mx       sync.RWMutex
	items    map[string]*item
	size     byte
	capacity byte
	order
}

func NewCache(cacheSize byte) (*Cache, error) {
	if cacheSize < 2 {
		return nil, errors.New("incorrect cache size")
	}

	return &Cache{items: make(map[string]*item, cacheSize), capacity: cacheSize}, nil
}

// Increment increments the expression counter if an item with such an index exists in the cache or returns false
func (c *Cache) Increment(index string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	var gottenItem *item
	gottenItem, ok = c.items[index]
	if ok {
		gottenItem.count++

		if c.order.head != gottenItem {

			// Set prev. and next fields for the items around
			if c.order.tail != gottenItem {
				gottenItem.previous.next, gottenItem.next.previous = gottenItem.next, gottenItem.previous
			} else {
				gottenItem.previous.next = nil
			}

			// Move the item to the beginning of the order
			gottenItem.previous = nil
			gottenItem.next = c.order.head

			c.order.head.previous = gottenItem
			c.order.head = gottenItem
		}
	}
	return
}

// Wipe removes every item of the cache saving the Cache size.
// The original items map will be garbage-collected eventually; you don't need to clear it manually.
func (c *Cache) Wipe() {
	*c = Cache{
		items:    make(map[string]*item, c.capacity),
		size:     0,
		capacity: c.capacity,
		order:    order{head: nil, tail: nil},
	}
}

// Delete deletes an item Cache with the index in the signature
func (c *Cache) Delete(index string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	var gottenItem *item
	gottenItem, ok = c.items[index]
	if !ok {
		return
	}

	if gottenItem.previous != nil {
		if gottenItem.next != nil {
			gottenItem.previous.next, gottenItem.next.previous = gottenItem.next, gottenItem.previous
		} else {
			gottenItem.previous.next = nil
			c.order.tail = gottenItem.previous
		}
	} else {
		if gottenItem.next == nil { // It means the item is the last
			c.order.head = nil
			c.order.tail = nil
		} else {
			gottenItem.next.previous = nil
			c.order.head = gottenItem.next
		}
	}

	gottenItem.previous = nil
	gottenItem.next = nil

	delete(c.items, index)
	c.size--
	return
}

// Add adds the new item to the Cache. Trows away the oldest item unless the Cache has free space.
// The data will be overwritten
func (c *Cache) Add(index string, data interface{}) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// New item creation
	item := item{count: 1, index: index, data: data}

	// Check if we have free space
	ok = c.capacity > c.size
	if ok {

		// Check if maximum count is reached
		if c.size != 255 {
			c.size++
		}
	} else {

		// Delete in the list the oldest item
		itemToDelete := c.order.tail
		delete(c.items, itemToDelete.index)

		// Delete the oldest item in the order
		itemToDelete.previous.next = nil
		c.order.tail = itemToDelete.previous
	}

	// add the new item to the cache
	c.items[index] = &item

	// add the new item to the order
	c.order.add(&item)

	return
}

// GetTheOldestIndex returns the oldest index in the cache
func (c *Cache) GetTheOldestIndex() string {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if c.order.tail != nil {
		return c.order.tail.index
	}

	return ""
}

// GetTheCacheSize returns the current Cache size that cannot be bigger than capacity
func (c *Cache) GetTheCacheSize() byte {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.size
}

// GetTheHeadIndex returns the current head index
func (c *Cache) GetTheHeadIndex() string {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if c.order.head != nil {
		return c.order.head.index
	}

	return ""
}

// GetTheNextItemIndex returns the index of the next item of the pointed out item's index
func (c *Cache) GetTheNextItemIndex(index string) (nextIndex string) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	item, ok := c.items[index]
	if ok && item.next != nil {
		nextIndex = item.next.index
	}
	return
}

// GetThePreviousItemIndex returns the index of the next item of the pointed out item's index
func (c *Cache) GetThePreviousItemIndex(index string) (previousIndex string) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	item, ok := c.items[index]
	if ok && item.previous != nil {
		previousIndex = item.previous.index
	}
	return
}

// GetStoredData() returns the related data stored in the Cache if the index is in the cache
func (c *Cache) GetStoredData(index string) (data interface{}, ok bool) {

	var item *item
	item, ok = c.items[index]
	if ok {
		data = item.data
	}

	return
}

// GetTheOldestCount returns the oldest index count
func (c *Cache) GetTheOldestCount() byte {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if c.order.tail != nil {
		return c.order.tail.count
	}

	return 0
}

// GetTheItemCount returns the index count
func (c *Cache) GetTheItemCount(index string) (byte, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	item, ok := c.items[index]
	if ok {
		return item.count, ok
	}

	return 0, ok
}

// add is an internal package method to keep order of the items
func (o *order) add(item *item) {
	switch o.head {
	case nil: // The order association list is empty
		o.head = item
		o.tail = item
	default:
		item.next = o.head
		o.head.previous = item
		o.head = item
	}
}
