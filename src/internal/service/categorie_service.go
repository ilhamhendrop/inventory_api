package service

import (
	"context"
	"database/sql"
	"errors"
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"time"

	"github.com/google/uuid"
)

type categorieService struct {
	categorieRepo model.CategorieRepository
}

func NewCategorie(categorieRepo model.CategorieRepository) model.CategorieService {
	return &categorieService{
		categorieRepo: categorieRepo,
	}
}

// Index implements [model.CategorieService].
func (c *categorieService) Index(ctx context.Context) (categories []dto.CategorieData, err error) {
	results, err := c.categorieRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		categories = append(categories, dto.CategorieData{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return categories, nil
}

// Search implements [model.CategorieService].
func (c *categorieService) Search(ctx context.Context, keyword string) (categories []dto.CategorieData, err error) {
	results, err := c.categorieRepo.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		categories = append(categories, dto.CategorieData{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return categories, nil
}

// Detail implements [model.CategorieService].
func (c *categorieService) Detail(ctx context.Context, id string) (categorie dto.CategorieData, err error) {
	persisted, err := c.categorieRepo.FindById(ctx, id)
	if err != nil {
		return categorie, err
	}

	if persisted.ID == "" {
		return categorie, errors.New("Data tidak ditemukan")
	}

	return dto.CategorieData{
		ID:          persisted.ID,
		Name:        persisted.Name,
		Description: persisted.Description,
	}, nil
}

// Create implements [model.CategorieService].
func (c *categorieService) Create(ctx context.Context, req dto.CategorieCreated) error {
	categorie := model.Categorie{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	return c.categorieRepo.Save(ctx, &categorie)
}

// Update implements [model.CategorieService].
func (c *categorieService) Update(ctx context.Context, req dto.CategorieUpdate) error {
	persisted, err := c.categorieRepo.FindById(ctx, req.ID)
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	persisted.Name = req.Name
	persisted.Description = req.Description
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	return c.categorieRepo.Update(ctx, &persisted)
}

// Delete implements [model.CategorieService].
func (c *categorieService) Delete(ctx context.Context, id string) error {
	persisted, err := c.categorieRepo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	return c.categorieRepo.Delete(ctx, id)
}
