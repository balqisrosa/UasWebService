package product

type Product struct {
	ProductId int `json:"product_id"`
	Name string `json:"name"`
	Price int `json:"price"`
	Stock int `json:"stock"`
	Description string `json:"description"`
}