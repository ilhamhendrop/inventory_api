package model

import (
	"context"
	"database/sql"
	"inventory-app/internal/dto"
)

type Categorie struct {
	ID          string       `db:"id"`
	Name        string       `db:"name"`
	Description string       `db:"description"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type CategorieRepository interface {
	FindAll(ctx context.Context) (categries []Categorie, err error)
	FindById(ctx context.Context, id string) (categorie Categorie, err error)
	FindByIds(ctx context.Context, ids []string) (categories []Categorie, err error)
	Search(ctx context.Context, keyword string) (categories []Categorie, err error)
	Save(ctx context.Context, ct *Categorie) error
	Update(ctx context.Context, ct *Categorie) error
	Delete(ctx context.Context, id string) error
}

type CategorieService interface {
	Index(ctx context.Context) (categories []dto.CategorieData, err error)
	Detail(ctx context.Context, id string) (categorie dto.CategorieData, err error)
	Search(ctx context.Context, keyword string) (categories []dto.CategorieData, err error)
	Create(ctx context.Context, req dto.CategorieCreated) error
	Update(ctx context.Context, req dto.CategorieUpdate) error
	Delete(ctx context.Context, id string) error
}
