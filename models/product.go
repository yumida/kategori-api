package models

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      *int   `json:"price"`
	CategoryID *int   `json:"category_id"`
	Stock      *int   `json:"stock"`
}
