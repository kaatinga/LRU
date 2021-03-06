package LRU

import (
	"reflect"
	"sync"
	"testing"

	"github.com/kaatinga/calc"
)

func TestNewCache(t *testing.T) {
	var cacheSize1 byte = 5
	var cacheSize2 byte = 0
	tests := []struct {
		name      string
		cacheSize byte
		want      *Cache
		wantErr   bool
	}{
		{"ok", cacheSize1, &Cache{
			mx:       sync.RWMutex{},
			items:    make(map[string]*item, cacheSize1),
			size:     0,
			capacity: cacheSize1,
			order:    order{},
		}, false},
		{"!ok", cacheSize2, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCache(tt.cacheSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddDelete(t *testing.T) {

	c, err := NewCache(3)
	if err != nil {
		t.Errorf("testCache was not created")
	}

	var TheOldestIndex = "1+1"

	tests := []struct {
		index  string
		wantOk bool
	}{
		{TheOldestIndex, true},
		{"1+2", true},
		{"1+3", true},
	}

	var result int64
	var gottenIndex string
	var gotOk bool
	for _, tt := range tests {
		t.Run(tt.index, func(t *testing.T) {

			result, err = calc.Calc(tt.index)
			if err != nil {
				t.Errorf("Calc package returned an error")
			}

			gotOk = c.Add(tt.index, result)
			if gotOk != tt.wantOk {
				t.Errorf("Increment() = %v, want %v", gotOk, tt.wantOk)
			}

			data, ok := c.GetStoredData(tt.index)
			if !ok {
				t.Errorf("Data are not stored! data: %v, want %v", data.(int64), result)
			}

			if data.(int64) != result {
				t.Errorf("Data are not correct! data: %v, want %v", data.(int64), result)
			}

			gottenIndex = c.GetTheOldestIndex()
			if gottenIndex != TheOldestIndex {
				t.Errorf("Cache last index = %v, want %v", gottenIndex, TheOldestIndex)
			}
		})
	}

	t.Run("checks after deleting 1+1", func(t *testing.T) {

		gotOk = c.Delete(TheOldestIndex)
		if gotOk != true {
			t.Errorf("Delete returned %v, want %v", gotOk, true)
		}

		gottenIndex = c.GetTheOldestIndex()
		if gottenIndex != "1+2" {
			t.Errorf("Cache last index = %v, want %v", gottenIndex, "1+2")
		}

		if c.tail.index != "1+2" {
			t.Errorf("The tail index is %v, want %v", c.tail.index, "1+2")
		}

		if c.head.index != "1+3" {
			t.Errorf("The head index is %v, want %v", c.head.index, "1+3")
		}

		if c.head.previous != nil {
			t.Errorf("The head previous is %v, want %s", c.head.previous.index, "nil")
		}

		if c.tail.next != nil {
			t.Errorf("The tail next is %v, want %s", c.head.next.index, "nil")
		}

		if c.head.next == nil {
			t.Errorf("The head.next is nil!")
		}

		if c.tail.previous == nil {
			t.Errorf("The tail.previous is nil!")
		}

		if c.head.next != c.tail {
			t.Errorf("The head.next (%s) is not tail (%s)!", c.head.next.index, c.tail.index)
		}

		if c.tail.previous != c.head {
			t.Errorf("The head (%s) is not tail.previous (%s)!", c.head.index, c.tail.previous.index)
		}

		if c.GetTheNextItemIndex(c.head.index) != c.tail.index {
			t.Errorf("The next index of the head (%s) is not tail (%s)!", c.GetTheNextItemIndex(c.head.index), c.tail.index)
		}

		if c.GetThePreviousItemIndex(c.tail.index) != c.head.index {
			t.Errorf("The previous index of the tail (%s) is not head (%s)!", c.GetThePreviousItemIndex(c.head.index), c.tail.index)
		}

		if c.GetTheCacheSize() != 2 {
			t.Errorf("The cache size (%v) is not correct, must be %v!", c.GetTheCacheSize(), 2)
		}
	})

	t.Run("wipe", func(t *testing.T) {
		c.Wipe()

		if c.size != 0 {
			t.Errorf("The cache size (%v) is not correct, must be %v!", c.GetTheCacheSize(), 0)
		}

		if c.capacity != 3 {
			t.Errorf("The cache capacity (%v) is not correct, must be %v!", c.capacity, 3)
		}

		if !reflect.DeepEqual(c.items, map[string]*item{}) {
			t.Errorf("The cache items must be nil, now is: %v", c.items)
		}

		if c.head != nil {
			t.Errorf("The head is not nil, now is: %v", c.head)
		}

		if c.tail != nil {
			t.Errorf("The tail is not nil, now is: %v", c.tail)
		}
	})
}
