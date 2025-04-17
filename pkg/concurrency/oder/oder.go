package oder

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
}

func main() {
	
}
