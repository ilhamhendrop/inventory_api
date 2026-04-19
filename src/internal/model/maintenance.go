package model

import (
	"context"
	"database/sql"
	"inventory-app/internal/dto"
	"time"
)

type Maintenance struct {
	ID          string       `db:"id"`
	ProductID   string       `db:"product_id"`
	Description string       `db:"description"`
	Status      string       `db:"status"`
	StartDate   time.Time    `db:"start_date"`
	EndDate     time.Time    `db:"end_date"`
	Stock       int          `db:"stock"`
	UserID      string       `db:"user_id"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type MaintenanceSearch struct {
	ProductID string
	UserID    string
}

type MaintenanceRepository interface {
	FindAll(ctx context.Context, ms MaintenanceSearch) (maintenances []Maintenance, err error)
	FindById(ctx context.Context, id string, ms MaintenanceSearch) (maintenance Maintenance, err error)
	FindByIds(ctx context.Context, ids []string) (maintenances []Maintenance, err error)
	Search(ctx context.Context, keyword string, ms MaintenanceSearch) (maintenance []Maintenance, err error)
	Save(ctx context.Context, mt *Maintenance) error
	Update(ctx context.Context, mt *Maintenance) error
	Delete(ctx context.Context, id string) error
}

type MaintenanceService interface {
	Index(ctx context.Context, ms MaintenanceSearch) (maintenances []dto.MaintenanceData, err error)
	Search(ctx context.Context, keyword string, ms MaintenanceSearch) (maintenances []dto.MaintenanceData, err error)
	Detail(ctx context.Context, id string, ms MaintenanceSearch) (maintenance dto.MaintenanceData, err error)
	Create(ctx context.Context, userId string, warehouseId string, req dto.MaintenanceCreated) error
	Update(ctx context.Context, req dto.MaintenanceUpdated) error
	Delete(ctx context.Context, id string) error
}
