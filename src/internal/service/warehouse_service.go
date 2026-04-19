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

type warehouseService struct {
	warehouseRepo model.WarehouseRepository
	productRepo   model.ProductRepository
	categorieRepo model.CategorieRepository
}

func NewWarehouse(
	warehouseRepo model.WarehouseRepository,
	productRepo model.ProductRepository,
	categorieRepo model.CategorieRepository) model.WarehouseService {
	return &warehouseService{
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
		categorieRepo: categorieRepo,
	}
}

// Index implements [model.WarehouseService].
func (w *warehouseService) Index(ctx context.Context, ws model.WarehouseSearch) (warehouses []dto.WarehouseData, err error) {
	results, err := w.warehouseRepo.FindAll(ctx, ws)
	if err != nil {
		return warehouses, err
	}

	productIds := make([]string, 0)
	for _, v := range results {
		productIds = append(productIds, v.ProductId)
	}

	products := make(map[string]model.Product)
	if len(productIds) > 0 {
		productData, err := w.productRepo.FindByIds(ctx, productIds)
		if err != nil {
			return warehouses, err
		}

		for _, v := range productData {
			products[v.ID] = v
		}
	}

	categorieIds := make([]string, 0)
	for _, v := range products {
		if v.CategorieId != "" {
			categorieIds = append(categorieIds, v.CategorieId)
		}
	}

	categories := make(map[string]model.Categorie)
	if len(categorieIds) > 0 {
		categorieData, err := w.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return warehouses, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	warehouseData := make([]dto.WarehouseData, 0)
	for _, v := range results {
		var product *dto.ProductData
		if vProduct, e := products[v.ProductId]; e {
			var categorie *dto.CategorieData
			if vCategorie, e := categories[vProduct.CategorieId]; e {
				categorie = &dto.CategorieData{
					ID:          vCategorie.ID,
					Name:        vCategorie.Name,
					Description: vCategorie.Description,
				}
			}

			product = &dto.ProductData{
				ID:        vProduct.ID,
				Name:      vProduct.Name,
				Categorie: categorie,
				Merek:     vProduct.Merek,
			}
		}

		warehouseData = append(warehouseData, dto.WarehouseData{
			ID:      v.ID,
			Product: product,
			Stock:   v.Stock,
			Status:  v.Status,
		})
	}

	return warehouseData, nil
}

// Search implements [model.WarehouseService].
func (w *warehouseService) Search(ctx context.Context, keyword string, ws model.WarehouseSearch) (warehouses []dto.WarehouseData, err error) {
	results, err := w.warehouseRepo.Search(ctx, keyword, ws)
	if err != nil {
		return warehouses, err
	}

	productIds := make([]string, 0)
	for _, v := range results {
		productIds = append(productIds, v.ProductId)
	}

	products := make(map[string]model.Product)
	if len(productIds) > 0 {
		productData, err := w.productRepo.FindByIds(ctx, productIds)
		if err != nil {
			return warehouses, err
		}

		for _, v := range productData {
			products[v.ID] = v
		}
	}

	categorieIds := make([]string, 0)
	for _, v := range products {
		if v.CategorieId != "" {
			categorieIds = append(categorieIds, v.CategorieId)
		}
	}

	categories := make(map[string]model.Categorie)
	if len(categorieIds) > 0 {
		categorieData, err := w.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return warehouses, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	warehouseData := make([]dto.WarehouseData, 0)
	for _, v := range results {
		var product *dto.ProductData
		if vProduct, e := products[v.ProductId]; e {
			var categorie *dto.CategorieData
			if vCategorie, e := categories[vProduct.CategorieId]; e {
				categorie = &dto.CategorieData{
					ID:          vCategorie.ID,
					Name:        vCategorie.Name,
					Description: vCategorie.Description,
				}
			}

			product = &dto.ProductData{
				ID:        vProduct.ID,
				Name:      vProduct.Name,
				Categorie: categorie,
				Merek:     vProduct.Merek,
			}
		}

		warehouseData = append(warehouseData, dto.WarehouseData{
			ID:      v.ID,
			Product: product,
			Stock:   v.Stock,
			Status:  v.Status,
		})
	}

	return warehouseData, nil
}

// Detail implements [model.WarehouseService].
func (w *warehouseService) Detail(ctx context.Context, id string, ws model.WarehouseSearch) (warehouse dto.WarehouseData, err error) {
	persisted, err := w.warehouseRepo.FindById(ctx, id, ws)
	if err != nil {
		return warehouse, err
	}

	if persisted.ID == "" {
		return warehouse, errors.New("Data tidak ditemukan")
	}

	productMdl, err := w.productRepo.FindById(ctx, persisted.ProductId, model.ProductSearch{})
	if err != nil {
		return warehouse, err
	}

	var categorie *dto.CategorieData
	if productMdl.CategorieId != "" {
		categories, err := w.categorieRepo.FindByIds(ctx, []string{productMdl.CategorieId})
		if err != nil {
			return warehouse, err
		}

		if len(categories) > 0 {
			categorie = &dto.CategorieData{
				ID:          categories[0].ID,
				Name:        categories[0].Name,
				Description: categories[0].Description,
			}
		}
	}

	product := &dto.ProductData{
		ID:        productMdl.ID,
		Name:      productMdl.Name,
		Categorie: categorie,
		Merek:     productMdl.Merek,
	}

	return dto.WarehouseData{
		ID:      persisted.ID,
		Product: product,
		Stock:   persisted.Stock,
		Status:  persisted.Status,
	}, nil
}

// Create implements [model.WarehouseService].
func (w *warehouseService) Create(ctx context.Context, req dto.WarehouseCreated) error {
	warehouse := model.Warehouse{
		ID:        uuid.NewString(),
		ProductId: req.ProductID,
		Stock:     req.Stock,
		Status:    model.WarehouseStatusReady,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	return w.warehouseRepo.Save(ctx, &warehouse)
}

// Update implements [model.WarehouseService].
func (w *warehouseService) Update(ctx context.Context, req dto.WarehouseUpdated) error {
	persisted, err := w.warehouseRepo.FindById(ctx, req.ID, model.WarehouseSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	persisted.ProductId = req.ProductID
	persisted.Stock = req.Stock
	persisted.Status = req.Status
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	return w.warehouseRepo.Update(ctx, &persisted)
}

// Delete implements [model.WarehouseService].
func (w *warehouseService) Delete(ctx context.Context, id string) error {
	persisted, err := w.warehouseRepo.FindById(ctx, id, model.WarehouseSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	return w.warehouseRepo.Delete(ctx, id)
}
