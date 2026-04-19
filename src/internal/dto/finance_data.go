package dto

type FinanceData struct {
	ID          string           `json:"id"`
	Maintenance *MaintenanceData `json:"maintenance"`
	Description string           `json:"description"`
	Price       int              `json:"price"`
	User        *UserData        `json:"user"`
}

type FinanceCreated struct {
	MaintenanceId string `json:"maintenance_id"`
	Description   string `json:"description" validate:"required"`
	Price         int    `json:"price" validate:"required"`
	UserId        string `json:"user_id"`
}

type FinanceUpdated struct {
	ID          string `json:"-"`
	Description string `json:"description" validate:"required"`
	Price       int    `json:"price" validate:"required"`
}
