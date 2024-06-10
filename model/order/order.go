package order

type Order struct {
	OrderId int `json:"order_id"`
	UserId int `json:"user_id"`
	Total float64 `json:"total"`
	Status string `json:"status"`
}