package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type StatusOrder string

const (
	StatusPending   StatusOrder = "pending"
	StatusShipped   StatusOrder = "shipped"
	StatusDelivered StatusOrder = "delivered"
	StatusCancelled StatusOrder = "cancelled"
)

type Order struct {
	ID     int
	Status StatusOrder
	mu     sync.Mutex
}

//var (
//	totalUpdates int
//	updateMutex  sync.Mutex
//)

func main() {

	var wg sync.WaitGroup
	wg.Add(2)

	orderChan := make(chan *Order)

	//orders := generateOrders(20)

	go func() {
		defer wg.Done()
		for _, order := range generateOrders(20) {
			orderChan <- order
		}
		close(orderChan)
	}()

	go processOrders(orderChan, &wg)

	//reportOrderStatus(orderChan)

	wg.Wait()
	fmt.Println("All operations completed. Exiting.")
	//fmt.Println("Total updates: ", totalUpdates)

}

func updateOrderStatuses(order *Order) {
	order.mu.Lock()
	statusOptions := []StatusOrder{
		StatusPending,
		StatusShipped,
		StatusDelivered,
		StatusCancelled,
	}
	time.Sleep(
		time.Duration(rand.Intn(300)) * time.Millisecond)
	newStatus := statusOptions[rand.Intn(len(statusOptions))]
	order.Status = newStatus
	fmt.Printf("Updated order %d status: %s\n", order.ID, newStatus)
	order.mu.Unlock()

	//updateMutex.Lock()
	//defer updateMutex.Unlock()
	//currentUpdates := totalUpdates
	//time.Sleep(5 * time.Millisecond)
	//totalUpdates = currentUpdates + 1
}

// orderChan <-chan *Order kenh chi nhan gia tri khong gui gia tri
// orderChan chan <- *Order kenh chi gui gia tri
// orderChan chan *Order kenh vua nhan va gui gia tri

func processOrders(orderChan <-chan *Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range orderChan {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		fmt.Printf("Processing orders %d\n", order.ID)
	}
}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{
			ID:     i + 1,
			Status: StatusPending,
		}
	}
	return orders
}

func reportOrderStatus(order []*Order) {
	for i := 0; i < 1; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("\n ------ Reported order status: ------ ")
		for _, order := range order {
			fmt.Printf("Order %d: %s\n", order.ID, order.Status)
		}
	}
	fmt.Println("-------------------------------------------------- ")
}
