package model

import (
	"context"
	"database/sql"
	"inventory-app/internal/dto"
)

type Product struct {
	ID          string       `db:"id"`
	Name        string       `db:"name"`
	CategorieId string       `db:"categorie_id"`
	Merek       string       `db:"merek"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdateAt    sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type ProductSearch struct {
	CategorieId string
}

type ProductRepository interface {
	FindAll(ctx context.Context, ps ProductSearch) (products []Product, err error)
	FindById(ctx context.Context, id string, ps ProductSearch) (product Product, err error)
	FindByIds(ctx context.Context, ids []string) (products []Product, err error)
	Search(ctx context.Context, keyword string, ps ProductSearch) (products []Product, err error)
	Save(ctx context.Context, pr *Product) error
	Update(ctx context.Context, pr *Product) error
	Delete(ctx context.Context, id string) error
}

type ProductService interface {
	Index(ctx context.Context, ps ProductSearch) (products []dto.ProductData, err error)
	Search(ctx context.Context, keyword string, ps ProductSearch) (products []dto.ProductData, err error)
	Detail(ctx context.Context, id string, ps ProductSearch) (product dto.ProductData, err error)
	Create(ctx context.Context, req dto.ProductCreated) error
	Update(ctx context.Context, req dto.ProductUpdated) error
	Delete(ctx context.Context, id string) error
}
