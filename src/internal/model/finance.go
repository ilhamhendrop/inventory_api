package model

import (
	"context"
	"database/sql"
	"inventory-app/internal/dto"
)

type Finance struct {
	ID            string       `db:"id"`
	MaintenanceId string       `db:"maintenance_id"`
	Description   string       `db:"description"`
	Price         int          `db:"price"`
	UserId        string       `db:"user_id"`
	CreatedAt     sql.NullTime `db:"created_at"`
	UpdatedAt     sql.NullTime `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
}

type FinanceSearch struct {
	MaintenanceId string
	UserId        string
}

type FinanceRepository interface {
	FindAll(ctx context.Context, fs FinanceSearch) (finances []Finance, err error)
	FindById(ctx context.Context, id string, fs FinanceSearch) (finance Finance, err error)
	Search(ctx context.Context, keyword string, fs FinanceSearch) (finances []Finance, err error)
	Save(ctx context.Context, fc *Finance) error
	Update(ctx context.Context, fc *Finance) error
	Delete(ctx context.Context, id string) error
}

type FinanceService interface {
	Index(ctx context.Context, fs FinanceSearch) (finances []dto.FinanceData, err error)
	Search(ctx context.Context, keyword string, fs FinanceSearch) (finances []dto.FinanceData, err error)
	Detail(ctx context.Context, id string, fs FinanceSearch) (finance dto.FinanceData, err error)
	Create(ctx context.Context, userId string, maintenanceId string, req dto.FinanceCreated) error
	Update(ctx context.Context, req dto.FinanceUpdated) error
	Delete(ctx context.Context, id string) error
}
