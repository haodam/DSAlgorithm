package main

import (
	"fmt"
	"runtime"
	"weak"
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

// Used weak pointers

type ResourceCacheWeak struct {
	items map[string]weak.Pointer[string]
}

func NewResourceCacheWeak() *ResourceCacheWeak {
	return &ResourceCacheWeak{
		items: make(map[string]weak.Pointer[string]),
	}
}

func (rc *ResourceCacheWeak) AddWeakPointers(key string, value *string) {
	weakPtr := weak.Make(value)
	rc.items[key] = weakPtr
}

func (rc *ResourceCacheWeak) GetWeakPointers(key string) (*string, bool) {
	weakPtr, exists := rc.items[key]
	if !exists {
		return nil, false
	}

	if ptr := weakPtr.Value(); ptr != nil {
		return ptr, true
	}
	delete(rc.items, key)
	return nil, false
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
	cache2 := NewResourceCacheWeak()

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

	cache2.AddWeakPointers("bigData", bigData)
	if _, exists := cache2.GetWeakPointers("bigData"); exists {
		fmt.Println("found big in cache 2")
	}

	bigData = nil
	fmt.Print("---------Raw Pointer---------:\n")
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

	fmt.Print("---------Weak Pointer---------:\n")
	fmt.Print("Before GC:")
	printMemoryUsage()

	runtime.GC()
	fmt.Print("After GC:")
	printMemoryUsage()

	if val, exists := cache2.GetWeakPointers("bigData"); exists {
		fmt.Printf(
			"Still holding %2.f MB in cache:\n",
			float64(len(*val)/1024/1024))
	}

}
