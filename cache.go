package main

import (
	"log"
	"sync"
	"time"
)

type Item struct {
	key string
	value int
	inserted int64
}

type Cache struct{
	Items []*Item
	mu  sync.RWMutex
}

// add the value to cache
func (c *Cache) Add(key string, value int){
	c.mu.Lock()
	data := &Item{}
	//e := time.Now().Add(expiry).UnixNano()
	e := time.Now().UnixNano()
	data.key = key
	data.value = value
	data.inserted = e
	cache.Items = append(cache.Items, data)
	log.Printf("item added to cache key %v, value is %v",key,value)
    c.mu.Unlock()
}


// returns array of values associated with key that are added in last duration
func (c *Cache) get(key string, duration time.Duration) ([]int,bool){
	c.mu.RLock()
	var values []int
	for _,item := range cache.Items{
		if(item.key == key && (item.inserted > 0)){

			if time.Now().Add(-duration).UnixNano() <= item.inserted {
				values = append(values,item.value)
	        }
		}
	}
	if len(values) > 0 {
		log.Printf("in getcache values are: %v", values)
		c.mu.RUnlock()
		return values,true
	} else {
		log.Printf("no values")
		c.mu.RUnlock()
		return nil,false
	}
}
