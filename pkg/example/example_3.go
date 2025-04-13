package main

import (
	"fmt"
	"sync"
)

func riskyNamedReturn() (result int) {
	var wg sync.WaitGroup
	result = 10

	localCopy := result //fix

	wg.Add(1)
	go func(val int) {
		defer wg.Done()
		fmt.Println("Goroutine reading result:", val) // 🔥 Đọc biến shared result
	}(localCopy)

	result = 20 // 🔥 Ghi đè lên cùng biến result
	wg.Wait()
	return
}

func Foo(id string) string {
	return id
}
func ProcessAll(uuids []string) {
	var myResults []string
	var mutex = &sync.Mutex{}
	safeAppend := func(res string) {
		mutex.Lock()
		myResults = append(myResults, res)
		mutex.Unlock()
	}
	for _, uuid := range uuids {
		go func(id string, results []string) {
			res := Foo(id) // Do something
			safeAppend(res)
		}(uuid, myResults) // su dung slice de truyen vao goroutine
	}
}

func ProcessAllFix(uuids []string) {
	var myResults []string
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	safeAppend := func(res string) {
		mutex.Lock()
		myResults = append(myResults, res)
		mutex.Unlock()
	}
	for _, uuid := range uuids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			res := Foo(id) // Do something
			safeAppend(res)
		}(uuid)
	}
}

func main() {
	fmt.Println("Final result:", riskyNamedReturn())
}
