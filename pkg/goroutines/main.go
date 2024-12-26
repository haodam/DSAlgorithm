package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	OrderId string
	Title   string
	Price   int
}

func Publisher(channel chan<- Message, orders []Message, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, order := range orders {
		fmt.Printf("Publishing%s\n", order.OrderId)
		channel <- order
		time.Sleep(1 * time.Second)
	}
	close(channel)
}

func subscriber(channel <-chan Message, userName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range channel {
		fmt.Printf("userName %s :: Order:: %s:: Price::%d\n", userName, msg.OrderId, msg.Price)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	orderChannel := make(chan Message)
	wg := sync.WaitGroup{}

	var orders = []Message{
		{OrderId: "Order-01", Title: "Hello World", Price: 1.0},
		{OrderId: "Order-02", Title: "Hello World", Price: 2.0},
		{OrderId: "Order-03", Title: "Hello World", Price: 3.0},
	}

	wg.Add(2)

	go Publisher(orderChannel, orders, &wg)
	go subscriber(orderChannel, "dam-anh-hao", &wg)
	time.Sleep(3 * time.Second)
	fmt.Println("Done")

}
