package Coffee

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
	Espresso CoffeesType = iota
	Cappuccino
	Latte
)

type Ingredient struct {
	Ingredient string // Ingredient: Nguyen lieu
	Quantity   int
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
	UpdateCoffeeQuantity(coffeeType CoffeesType, quantity int) error
	GetCoffee(coffeeType CoffeesType) (*Coffee, error)
	GetAllCoffeeAvailable() []Coffee
}

type InMemoryCoffeeRepository struct {
	coffees map[CoffeesType]*Coffee
	//observers []Observer
}

func main() {

}
