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

type productService struct {
	productRepo   model.ProductRepository
	categorieRepo model.CategorieRepository
}

func NewProduct(productRepo model.ProductRepository, categorieRepo model.CategorieRepository) model.ProductService {
	return &productService{
		productRepo:   productRepo,
		categorieRepo: categorieRepo,
	}
}

// Index implements [model.ProductService].
func (p *productService) Index(ctx context.Context, ps model.ProductSearch) (products []dto.ProductData, err error) {
	results, err := p.productRepo.FindAll(ctx, ps)
	if err != nil {
		return products, err
	}

	categorieIds := make([]string, 0)
	for _, v := range results {
		categorieIds = append(categorieIds, v.CategorieId)
	}

	categories := make(map[string]model.Categorie)
	if len(categorieIds) > 0 {
		categorieData, err := p.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return products, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	productData := make([]dto.ProductData, 0)
	for _, v := range results {
		var categorie *dto.CategorieData
		if vCategorie, e := categories[v.CategorieId]; e {
			categorie = &dto.CategorieData{
				ID:          vCategorie.ID,
				Name:        vCategorie.Name,
				Description: vCategorie.Description,
			}
		}

		productData = append(productData, dto.ProductData{
			ID:        v.ID,
			Name:      v.Name,
			Categorie: categorie,
			Merek:     v.Merek,
		})
	}

	return productData, nil
}

// Search implements [model.ProductService].
func (p *productService) Search(ctx context.Context, keyword string, ps model.ProductSearch) (products []dto.ProductData, err error) {
	results, err := p.productRepo.Search(ctx, keyword, ps)
	if err != nil {
		return products, err
	}

	categorieIds := make([]string, 0)
	for _, v := range results {
		categorieIds = append(categorieIds, v.CategorieId)
	}

	categories := make(map[string]model.Categorie)
	if len(categorieIds) > 0 {
		categorieData, err := p.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return products, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	productData := make([]dto.ProductData, 0)
	for _, v := range results {
		var categorie *dto.CategorieData
		if vCategorie, e := categories[v.CategorieId]; e {
			categorie = &dto.CategorieData{
				ID:          vCategorie.ID,
				Name:        vCategorie.Name,
				Description: vCategorie.Description,
			}
		}

		productData = append(productData, dto.ProductData{
			ID:        v.ID,
			Name:      v.Name,
			Categorie: categorie,
			Merek:     v.Merek,
		})
	}

	return productData, nil
}

// Detail implements [model.ProductService].
func (p *productService) Detail(ctx context.Context, id string, ps model.ProductSearch) (product dto.ProductData, err error) {
	persisted, err := p.productRepo.FindById(ctx, id, ps)
	if err != nil {
		return product, err
	}

	if persisted.ID == "" {
		return product, errors.New("Data categorie tidak ditemukan")
	}

	var categorie *dto.CategorieData
	if persisted.CategorieId != "" {
		categories, err := p.categorieRepo.FindByIds(ctx, []string{persisted.CategorieId})
		if err != nil {
			return product, err
		}

		if len(categories) > 0 {
			categorie = &dto.CategorieData{
				ID:          categories[0].ID,
				Name:        categories[0].Name,
				Description: categories[0].Description,
			}
		}
	}

	return dto.ProductData{
		ID:        persisted.ID,
		Name:      product.Name,
		Categorie: categorie,
		Merek:     persisted.Merek,
	}, nil
}

// Create implements [model.ProductService].
func (p *productService) Create(ctx context.Context, req dto.ProductCreated) error {
	product := model.Product{
		ID:          uuid.NewString(),
		Name:        req.Name,
		CategorieId: req.CategorieId,
		Merek:       req.Merek,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	return p.productRepo.Save(ctx, &product)
}

// Update implements [model.ProductService].
func (p *productService) Update(ctx context.Context, req dto.ProductUpdated) error {
	persisted, err := p.productRepo.FindById(ctx, req.ID, model.ProductSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	persisted.Name = req.Name
	persisted.CategorieId = req.CategorieId
	persisted.Merek = req.Merek
	persisted.UpdateAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	return p.productRepo.Update(ctx, &persisted)
}

// Delete implements [model.ProductService].
func (p *productService) Delete(ctx context.Context, id string) error {
	persisted, err := p.productRepo.FindById(ctx, id, model.ProductSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	return p.productRepo.Delete(ctx, id)
}
