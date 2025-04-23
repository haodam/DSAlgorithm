package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func placeOrderWithoutContext(ctx context.Context, orderID string) error {
	log.Printf("bắt đầu xử lý đơn hàng: %s\n", orderID)
	select {
	case <-ctx.Done():
		log.Printf("xử lý đơn hàng %s: %v\n", orderID, ctx.Err())
		return ctx.Err()
	case <-time.After(time.Second * 3):
		log.Printf("xử lý đơn hàng %s thành công \n", orderID)
		return nil
	}
}

func OrderHandlerWithContext(w http.ResponseWriter, r *http.Request) {
	orderID := "Order-12345"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := placeOrderWithoutContext(ctx, orderID)
	if err != nil {
		log.Printf("xử lý đơn hàng %s thất bại: %s\n", orderID, err)
		http.Error(w, "xử lý đơn hàng hoặc quá thời gian", http.StatusGatewayTimeout)
		return
	}
	w.WriteHeader(http.StatusOK)
	write, err := w.Write([]byte("ĐẶT HÀNG THÀNH CÔNG" + orderID))
	if err != nil {
		return
	}
	fmt.Println(write)
}

func main() {
	http.HandleFunc("/order", OrderHandlerWithContext)
	log.Println("Server starting... http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
