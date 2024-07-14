package order

type Order struct {
	OrderId int `json:"order_id"`
	UserId int `json:"user_id"`
	ProductId int `json:"product_id"`
	Amount int `json:"amount"`
	Price float64 `json:"price"`
	Total float64 `json:"total"`
	Status string `json:"status"`
}