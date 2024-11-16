package main

import (
	"fmt"
	"sync"
)

// Designing a Coffee Vending Machine

// The coffee vending machine should support different types of coffee, such as espresso, cappuccino, and latte.
// Each type of coffee should have a specific price and recipe (ingredients and their quantities).
// The machine should have a menu to display the available coffee options and their prices.
// Users should be able to select a coffee type and make a payment.
// The machine should dispense the selected coffee and provide change if necessary.
// The machine should track the inventory of ingredients and notify when they are running low.
// The machine should handle multiple user requests concurrently and ensure thread safety.

// Thiết kế máy bán cà phê tự động

// Máy bán cà phê tự động phải hỗ trợ nhiều loại cà phê khác nhau, chẳng hạn như espresso, cappuccino và latte.
// Mỗi loại cà phê phải có giá và công thức cụ thể (thành phần và số lượng).
// Máy phải có menu để hiển thị các loại cà phê có sẵn và giá của chúng.
// Người dùng phải có thể chọn loại cà phê và thực hiện thanh toán.
// Máy phải phân phối cà phê đã chọn và trả lại tiền thừa nếu cần.
// Máy phải theo dõi lượng nguyên liệu tồn kho và thông báo khi chúng sắp hết.
// Máy phải xử lý nhiều yêu cầu của người dùng đồng thời và đảm bảo an toàn luồng.

type CoffeesType int

const (
	Espresso CoffeesType = iota + 1
	Cappuccino
	Latte
)

type Ingredient struct {
	Name     string // Ingredient: Nguyen lieu
	Quantity int
}

type Coffee struct {
	Name     string
	Type     CoffeesType
	Price    int
	Recipe   []Ingredient // Recipe: cong thuc
	Quantity int
}

type CoffeesRepository interface {
	AddCoffee(coffee Coffee) error
	GetCoffee(coffeeType CoffeesType) (*Coffee, error)
	GetAllCoffeeAvailable() []Coffee
}

// Lưu giữ dữ liệu vào trong memory thay cho csdl

type InMemoryCoffeeRepository struct {
	coffees map[CoffeesType]*Coffee
}

func (repo *InMemoryCoffeeRepository) AddCoffee(coffee Coffee) error {

	if coffee.Price <= 0 {
		return fmt.Errorf("coffee price must be positive")
	}

	_, found := repo.coffees[coffee.Type]
	if !found {
		repo.coffees[coffee.Type] = &coffee
	}
	return nil
}

func (repo *InMemoryCoffeeRepository) GetCoffee(coffeeType CoffeesType) (*Coffee, error) {

	coffee, found := repo.coffees[coffeeType]
	if !found {
		return nil, fmt.Errorf("coffee with type %s does not exist", coffeeType)
	}
	return coffee, nil
}

func (repo *InMemoryCoffeeRepository) GetAllCoffeeAvailable() []Coffee {

	var rs []Coffee
	for _, coffee := range repo.coffees {
		rs = append(rs, *coffee)
	}
	return rs
}

type Observer interface {
	Update(ingredient Ingredient)
}

type AdminObserver struct {
}

func (admin *AdminObserver) Update(ingredient Ingredient) {
	fmt.Println("\n You have notification")
	fmt.Printf("Ingredient %s is running low: %d \n", ingredient.Name, ingredient.Quantity)
}

type Coin struct {
	Value int
	Count int
}

type Payment struct {
	coins map[int]int
	mu    *sync.Mutex
}

type IPayment interface {
	Payment(coins []Coin) (int, error)
	ReturnChange(amount int) ([]Coin, error)
}

func (p *Payment) Payment(coins []Coin) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Sao luu trang thai hien tai (phong truong hop khoi phuc)
	backup := make(map[int]int)
	for k, v := range p.coins {
		backup[k] = v
	}

	total := 0
	for _, coin := range coins {
		if coin.Count <= 0 || coin.Value <= 0 {
			return coin.Value, fmt.Errorf("invalid coin amount")
		}
		p.coins[coin.Value] += coin.Count
		total += coin.Count

	}
	return total, nil
}

func main() {
	// Khởi tạo repository
	repo := &InMemoryCoffeeRepository{
		coffees: make(map[CoffeesType]*Coffee),
	}

	// Thêm loại cà phê hợp lệ
	err := repo.AddCoffee(Coffee{
		Name:  "Espresso",
		Type:  Espresso,
		Price: 50000,
	})
	if err != nil {
		fmt.Println("Lỗi:", err)
	} else {
		fmt.Println("Thêm cà phê thành công!")
	}

	// Thêm cà phê với giá không hợp lệ
	err = repo.AddCoffee(Coffee{
		Name:  "Invalid Coffee",
		Type:  Latte,
		Price: -10000,
	})
	if err != nil {
		fmt.Println("Lỗi:", err)
	}

	// Thêm cà phê trùng lặp
	err = repo.AddCoffee(Coffee{
		Name:  "Espresso",
		Type:  Espresso,
		Price: 50000,
	})
	if err != nil {
		fmt.Println("Lỗi:", err)
	} else {
		fmt.Println("Thêm cà phê trùng lặp vẫn không lỗi!")
	}

	// In danh sách cà phê trong repository
	for coffeeType, coffee := range repo.coffees {
		fmt.Printf("Loại: %v, Thông tin: %+v\n", coffeeType, *coffee)
	}

}
