package dto

type WarehouseData struct {
	ID      string       `json:"id"`
	Product *ProductData `json:"product"`
	Stock   int          `json:"stock"`
	Status  string       `json:"status"`
}

type WarehouseCreated struct {
	ProductID string `json:"product_id" validate:"required"`
	Stock     int    `json:"stock" validate:"required"`
	Status    string `json:"status"`
}

type WarehouseUpdated struct {
	ID        string `json:"-"`
	ProductID string `json:"product_id" validate:"required"`
	Stock     int    `json:"stock" validate:"required"`
	Status    string `json:"status" validate:"required"`
}
