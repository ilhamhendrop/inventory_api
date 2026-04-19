package dto

type MaintenanceData struct {
	ID          string       `json:"id"`
	Product     *ProductData `json:"product"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	Stock       int          `json:"stock"`
	User        *UserData    `json:"user"`
}

type MaintenanceCreated struct {
	ProductId   string `json:"product_id"`
	UserId      string `json:"user_id"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status"`
	StartDate   string `json:"start_date" validate:"required"`
	Stock       int    `json:"stock" validate:"required"`
}

type MaintenanceUpdated struct {
	ID          string `json:"-"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
	Stock       int    `json:"stock" validate:"required"`
}
