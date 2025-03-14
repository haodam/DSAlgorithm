package main

import (
	"fmt"
	"github.com/haodam/DSAlgorithm/pkg/polymorphism"
)

type purchasable interface {
	CalculatePrice() int64
}

var cart []purchasable

func addToCart(products ...purchasable) {
	cart = append(cart, products...)
}

func getCartTotal() int64 {
	var total = int64(0)
	for _, product := range cart {
		total += product.CalculatePrice()
	}
	return total
}
func main() {
	myShirt := polymorphism.Shirt{polymorphism.ProductDetails{Price: 5000, Brand: "Nike"}, "XL", "Red"}
	myMonitor := polymorphism.Monitor{polymorphism.ProductDetails{Price: 10000, Brand: "SamSung"}, "32 inch", "4k"}
	myWine := polymorphism.Wine{polymorphism.ProductDetails{Price: 90, Brand: "Tao Meo"}, "2000", "Red"}

	addToCart(myShirt, myMonitor, myWine)
	fmt.Println(getCartTotal())

}
