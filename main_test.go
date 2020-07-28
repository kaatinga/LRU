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

func TestCache_Add(t *testing.T) {

	c, err := NewCache(3)
	if err != nil {
		t.Errorf("testCache was not created")
	}

	var TheOldestIndex = "1+1"

	tests := []struct {
		index  string
		wantOk bool
	}{
		{ TheOldestIndex, true},
		{ "1+2", true},
		{ "1+3", true},
	}
	var result int64
	var gottenIndex string
		for _, tt := range tests {
		t.Run(tt.index, func(t *testing.T) {

			result, err = calc.Calc(tt.index)
			if err != nil {
				t.Errorf("Calc package returned an error")
			}

			if gotOk := c.Add(tt.index, result); gotOk != tt.wantOk {
				t.Errorf("Increment() = %v, want %v", gotOk, tt.wantOk)
			}

			gottenIndex = c.GetTheOldestIndex()
			if gottenIndex != TheOldestIndex {
				t.Errorf("Cache last index = %v, want %v", gottenIndex, TheOldestIndex)
			}
		})
	}
}
