package main

import (
	"fmt"
	"sync"
)

// ...get instance...
// Single get instance address: 0x21b4c0
// Single get instance address: 0x21b4c0
// 2 Single duoc xuat ra nhung 2 dia chi deu khac nhau

// Sau khi fix error , tra ve 1 duy nhat get instance va dia chi giong nhau
// ...get instance...
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0
// Singleton get instance address: 0x5bc4c0

type Singleton struct {
}

var (
	instance *Singleton
	once     sync.Once
)

func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Printf("...get instance...\n")
		instance = &Singleton{}
	})
	return instance
}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := GetInstance()
			fmt.Printf("Singleton get instance address: %p\n", s)
		}()
	}
	wg.Wait()
}
