package main

import (
	"fmt"
	"sync"
)

//func processJob(job string) {
//	fmt.Println(job)
//}
//
//func Foo() {}
//
//func ProcessAll(uuids []string) {
//	var myResults []string
//	var mutex sync.Mutex
//	safeAppend := func(res string) {
//		mutex.Lock()
//		myResults = append(myResults, res)
//		mutex.Unlock()
//	}
//
//	for _, uuid := range uuids {
//		go func(id string, results []string) {
//			res := Foo(id)
//			safeAppend(res)
//		}(uuid, myResults) // slice read without holding lock
//	}
//}

type Counter struct {
	mu  sync.Mutex
	val int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val++
	fmt.Println("value", c.val) // Sử dụng fmt.Println để in ra màn hình
}

//func processOrders(uuids []string)  {
//	var errMap = make(map[string]error)
//	for _, uuid := range uuids {
//		go func(uuid string) {
//		orderHandle, err := GetOrder(uuid)
//		if err != nil {
//			errMap[uuid] = err
//			return
//		}
//
//		...
//	}(uuid)
//		return combineErrors(errMap)
//	}
//}

func main() {

	var wg sync.WaitGroup
	counter := &Counter{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	//fmt.Sprintln("Final counter;", counter.val)

	//jobs := []string{"job1", "job2", "job3"}
	//
	//for _, job := range jobs {
	//	go func() {
	//		processJob(job)
	//	}()
	//}
	//time.Sleep(1 * time.Second)

	//x, err := Foo()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//go func() {
	//	var y int
	//
	//}()
}
