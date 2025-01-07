package event

type Trader struct {
	Pair  Pair
	Price float64
	Qty   float64
	IsBuy bool
	Unix  int64
}

type Pair struct {
	Exchange string
	Symbol   string
}
