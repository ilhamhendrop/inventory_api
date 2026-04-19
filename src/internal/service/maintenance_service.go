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

type maintenanceService struct {
	maintenanceRepo model.MaintenanceRepository
	productRepo     model.ProductRepository
	categorieRepo   model.CategorieRepository
	warehouseRepo   model.WarehouseRepository
	userRepo        model.UserRepository
}

func NewMaintenance(
	maintenanceRepo model.MaintenanceRepository,
	productRepo model.ProductRepository,
	categorieRepo model.CategorieRepository,
	warehouseRepo model.WarehouseRepository,
	userRepo model.UserRepository,
) model.MaintenanceService {
	return &maintenanceService{
		maintenanceRepo: maintenanceRepo,
		productRepo:     productRepo,
		categorieRepo:   categorieRepo,
		warehouseRepo:   warehouseRepo,
		userRepo:        userRepo,
	}
}

// Index implements [model.MaintenanceService].
func (m *maintenanceService) Index(ctx context.Context, ms model.MaintenanceSearch) (maintenances []dto.MaintenanceData, err error) {
	results, err := m.maintenanceRepo.FindAll(ctx, ms)
	if err != nil {
		return maintenances, err
	}

	productIds := make([]string, 0)
	userIds := make([]string, 0)
	for _, v := range results {
		productIds = append(productIds, v.ProductID)
		userIds = append(userIds, v.UserID)
	}

	products := make(map[string]model.Product)
	if len(productIds) > 0 {
		productData, err := m.productRepo.FindByIds(ctx, productIds)
		if err != nil {
			return maintenances, err
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
		categorieData, err := m.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return maintenances, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	users := make(map[string]model.User)
	if len(userIds) > 0 {
		userData, err := m.userRepo.FindByIds(ctx, userIds)
		if err != nil {
			return maintenances, err
		}

		for _, v := range userData {
			users[v.ID] = v
		}
	}

	maintenanceData := make([]dto.MaintenanceData, 0)
	for _, v := range results {
		var product *dto.ProductData
		if vProduct, e := products[v.ProductID]; e {
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

		var user *dto.UserData
		if vUser, e := users[v.UserID]; e {
			user = &dto.UserData{
				ID:       vUser.ID,
				Username: vUser.Username,
				Name:     vUser.Name,
				Role:     vUser.Role,
				Status:   vUser.Status,
			}
		}

		maintenanceData = append(maintenanceData, dto.MaintenanceData{
			ID:          v.ID,
			Product:     product,
			Description: v.Description,
			Status:      v.Status,
			StartDate:   v.StartDate.Format("2006-01-02"),
			EndDate:     v.EndDate.Format("2006-01-02"),
			Stock:       v.Stock,
			User:        user,
		})
	}

	return maintenanceData, nil
}

// Search implements [model.MaintenanceService].
func (m *maintenanceService) Search(ctx context.Context, keyword string, ms model.MaintenanceSearch) (maintenances []dto.MaintenanceData, err error) {
	results, err := m.maintenanceRepo.Search(ctx, keyword, ms)
	if err != nil {
		return maintenances, err
	}

	productIds := make([]string, 0)
	userIds := make([]string, 0)
	for _, v := range results {
		productIds = append(productIds, v.ProductID)
		userIds = append(userIds, v.UserID)
	}

	products := make(map[string]model.Product)
	if len(productIds) > 0 {
		productData, err := m.productRepo.FindByIds(ctx, productIds)
		if err != nil {
			return maintenances, err
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
		categorieData, err := m.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return maintenances, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	users := make(map[string]model.User)
	if len(userIds) > 0 {
		userData, err := m.userRepo.FindByIds(ctx, userIds)
		if err != nil {
			return maintenances, err
		}

		for _, v := range userData {
			users[v.ID] = v
		}
	}

	maintenanceData := make([]dto.MaintenanceData, 0)
	for _, v := range results {
		var product *dto.ProductData
		if vProduct, e := products[v.ProductID]; e {
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

		var user *dto.UserData
		if vUser, e := users[v.UserID]; e {
			user = &dto.UserData{
				ID:       vUser.ID,
				Username: vUser.Username,
				Name:     vUser.Name,
				Role:     vUser.Role,
				Status:   vUser.Status,
			}
		}

		maintenanceData = append(maintenanceData, dto.MaintenanceData{
			ID:          v.ID,
			Product:     product,
			Description: v.Description,
			Status:      v.Status,
			StartDate:   v.StartDate.Format("2006-01-02"),
			EndDate:     v.EndDate.Format("2006-01-02"),
			Stock:       v.Stock,
			User:        user,
		})
	}

	return maintenanceData, nil
}

// Detail implements [model.MaintenanceService].
func (m *maintenanceService) Detail(ctx context.Context, id string, ms model.MaintenanceSearch) (maintenance dto.MaintenanceData, err error) {
	persisted, err := m.maintenanceRepo.FindById(ctx, id, ms)
	if err != nil {
		return maintenance, err
	}

	if persisted.ID == "" {
		return maintenance, errors.New("Data tidak ditemukan")
	}

	productMd, err := m.productRepo.FindById(ctx, persisted.ProductID, model.ProductSearch{})
	if err != nil {
		return maintenance, err
	}

	userMd, err := m.userRepo.FindById(ctx, persisted.UserID)

	var categorie *dto.CategorieData
	if productMd.CategorieId != "" {
		categories, err := m.categorieRepo.FindByIds(ctx, []string{productMd.CategorieId})
		if err != nil {
			return maintenance, err
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
		ID:        productMd.ID,
		Name:      productMd.Name,
		Categorie: categorie,
		Merek:     productMd.Merek,
	}

	user := &dto.UserData{
		ID:       userMd.ID,
		Username: userMd.Username,
		Name:     userMd.Name,
		Role:     userMd.Role,
		Status:   userMd.Status,
	}

	return dto.MaintenanceData{
		ID:          persisted.ID,
		Product:     product,
		Description: persisted.Description,
		Status:      persisted.Status,
		StartDate:   persisted.StartDate.Format("2006-01-02"),
		EndDate:     persisted.EndDate.Format("2006-01-02"),
		Stock:       persisted.Stock,
		User:        user,
	}, nil
}

// Create implements [model.MaintenanceService].
func (m *maintenanceService) Create(ctx context.Context, userId string, warehouseId string, req dto.MaintenanceCreated) error {
	persisted, err := m.warehouseRepo.FindById(ctx, warehouseId, model.WarehouseSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data warehouse tidak ditemukan")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return err
	}

	maintenance := model.Maintenance{
		ID:          uuid.NewString(),
		ProductID:   persisted.ProductId,
		Description: req.Description,
		Status:      model.MaintenanceStatusDiproses,
		StartDate:   startDate,
		Stock:       req.Stock,
		UserID:      userId,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	return m.maintenanceRepo.Save(ctx, &maintenance)
}

// Update implements [model.MaintenanceService].
func (m *maintenanceService) Update(ctx context.Context, req dto.MaintenanceUpdated) error {
	persisted, err := m.maintenanceRepo.FindById(ctx, req.ID, model.MaintenanceSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return err
	}

	oldStatus := persisted.Status

	persisted.Description = req.Description
	persisted.Status = req.Status
	persisted.EndDate = endDate
	persisted.Stock = req.Stock
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	if err := m.maintenanceRepo.Update(ctx, &persisted); err != nil {
		return err
	}

	if oldStatus != model.MaintenanceStatusRusak && req.Status == model.MaintenanceStatusRusak {
		persistedWarehouse, err := m.warehouseRepo.FindByProductId(ctx, persisted.ProductID)
		if err != nil {
			return err
		}

		if persistedWarehouse.ID == "" {
			return errors.New("Data product di warehouse tidak ditemukan")
		}

		if req.Stock > persistedWarehouse.Stock {
			return errors.New("stock warehouse tidak mencukupi")
		}

		persistedWarehouse.Stock -= req.Stock
		persistedWarehouse.UpdatedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}

		if err := m.warehouseRepo.Update(ctx, &persistedWarehouse); err != nil {
			return err
		}
	}

	return nil
}

// Delete implements [model.MaintenanceService].
func (m *maintenanceService) Delete(ctx context.Context, id string) error {
	persisted, err := m.maintenanceRepo.FindById(ctx, id, model.MaintenanceSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	return m.maintenanceRepo.Delete(ctx, id)
}
