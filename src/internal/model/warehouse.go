package model

import (
	"context"
	"database/sql"
	"inventory-app/internal/dto"
)

type Warehouse struct {
	ID        string       `db:"id"`
	ProductId string       `db:"product_id"`
	Stock     int          `db:"stock"`
	Status    string       `db:"status"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type WarehouseSearch struct {
	ProductId string
}

type WarehouseRepository interface {
	FindAll(ctx context.Context, ws WarehouseSearch) (warehouses []Warehouse, err error)
	FindById(ctx context.Context, id string, ws WarehouseSearch) (warehouse Warehouse, err error)
	FindByIds(ctx context.Context, ids []string) (warehouses []Warehouse, err error)
	FindByProductId(ctx context.Context, productId string) (warehouse Warehouse, err error)
	Search(ctx context.Context, keyword string, ws WarehouseSearch) (warehouses []Warehouse, err error)
	Save(ctx context.Context, wh *Warehouse) error
	Update(ctx context.Context, wh *Warehouse) error
	Delete(ctx context.Context, id string) error
}

type WarehouseService interface {
	Index(ctx context.Context, ws WarehouseSearch) (warehouses []dto.WarehouseData, err error)
	Search(ctx context.Context, keyword string, ws WarehouseSearch) (warehouses []dto.WarehouseData, err error)
	Detail(ctx context.Context, id string, ws WarehouseSearch) (warehouse dto.WarehouseData, err error)
	Create(ctx context.Context, req dto.WarehouseCreated) error
	Update(ctx context.Context, req dto.WarehouseUpdated) error
	Delete(ctx context.Context, id string) error
}
