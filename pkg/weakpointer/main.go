package main

import (
	"fmt"
	"runtime"
)

type ResourceCache struct {
	items map[string]*string
}

func NewResourceCache() *ResourceCache {
	return &ResourceCache{
		items: make(map[string]*string),
	}
}

func (rc *ResourceCache) Add(key string, value *string) {
	rc.items[key] = value
}

func (rc *ResourceCache) Get(key string) (*string, bool) {
	value, exists := rc.items[key]
	return value, exists
}

func printMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf(
		"Memory: %.2f MB\n",
		float64(m.Alloc/1024/1024))
}

func main() {

	cache := NewResourceCache()

	createBigString := func() *string {
		s := make([]byte, 10<<20)
		str := string(s)
		return &str
	}

	bigData := createBigString()
	cache.Add("bigData", bigData)
	if _, exists := cache.Get("bigData"); exists {
		fmt.Println("found big in cache")
	}

	bigData = nil

	fmt.Print("Before GC:")
	printMemoryUsage()

	runtime.GC()
	fmt.Print("After GC:")
	printMemoryUsage()

	if val, exists := cache.Get("bigData"); exists {
		fmt.Printf(
			"Still holding %2.f MB in cache:\n",
			float64(len(*val)/1024/1024))
	}
}
