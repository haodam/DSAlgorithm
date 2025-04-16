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
		}(uuid, myResults) // 🔥 su dung slice de truyen vao goroutine
	}
}

// data race with slice

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
	wg.Wait()
}

func GetOrder(uuid string) ([]string, error) {
	return nil, fmt.Errorf("fake error for %s", uuid)
}

// data race with map
func processOrders(uuids []string) error {
	var errMap = make(map[string]error)
	for _, uuid := range uuids {
		go func(uuid string) {
			_, err := GetOrder(uuid)
			if err != nil {
				errMap[uuid] = err // 🔥
				return
			}
		}(uuid)
	}
	return errMap[""]
}

// fix data race with map
func processOrdersFix(uuids []string) error {
	var (
		errMap = make(map[string]error)
		mutex  = &sync.Mutex{}
		wg     sync.WaitGroup
	)
	for _, uuid := range uuids {
		wg.Add(1)
		go func(uuid string) {
			defer wg.Done()
			_, err := GetOrder(uuid)
			if err != nil {
				mutex.Lock()
				errMap[uuid] = err
				mutex.Unlock()
				return
			}
		}(uuid)
	}
	return errMap[""]
}

// data race with pass-by-value
var a int

// CriticalSection receives a copy of mutex
func CriticalSection(m *sync.Mutex) {
	m.Lock()
	a = a + 1
	m.Unlock()
}

// Example data race with pass-by-value

type Counter struct {
	mu sync.Mutex
	n  int
}

func Increment(c *Counter) {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func main() {
	//fmt.Println("Final result:", riskyNamedReturn())
	mutex := &sync.Mutex{}
	go CriticalSection(mutex)
	go CriticalSection(mutex)
}
