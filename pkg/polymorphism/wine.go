package polymorphism

type Wine struct {
	ProductDetails
	Year string
	Kind string // Loai ruou vang
}

func (w Wine) CalculatePrice() int64 {
	liquorTax := float64(w.Price) * .23
	stateLiquorTax := float64(w.Price) * .10
	return w.Price + int64(liquorTax) + int64(stateLiquorTax)
}
