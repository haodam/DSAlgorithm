package polymorphism

type Monitor struct {
	ProductDetails
	Size       string
	Resolution string // Do phan giai
}

func (m Monitor) CalculatePrice() int64 {
	electronicsTax := float64(m.Price) * .30
	return m.Price + int64(electronicsTax)
}
