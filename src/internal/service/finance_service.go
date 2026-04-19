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

type financeService struct {
	financeRepo     model.FinanceRepository
	maintenanceRepo model.MaintenanceRepository
	userRepo        model.UserRepository
	productRepo     model.ProductRepository
	categorieRepo   model.CategorieRepository
}

func NewFinance(
	financeRepo model.FinanceRepository,
	maintenanceRepo model.MaintenanceRepository,
	userRepo model.UserRepository,
	productRepo model.ProductRepository,
	categorieRepo model.CategorieRepository) model.FinanceService {
	return &financeService{
		financeRepo:     financeRepo,
		maintenanceRepo: maintenanceRepo,
		userRepo:        userRepo,
		productRepo:     productRepo,
		categorieRepo:   categorieRepo,
	}
}

// Index implements [model.FinanceService].
func (f *financeService) Index(ctx context.Context, fs model.FinanceSearch) (finances []dto.FinanceData, err error) {
	results, err := f.financeRepo.FindAll(ctx, fs)
	if err != nil {
		return finances, err
	}

	maintenanceIds := make([]string, 0)
	userIds := make([]string, 0)
	for _, v := range results {
		maintenanceIds = append(maintenanceIds, v.MaintenanceId)
		userIds = append(userIds, v.UserId)
	}

	maintenances := make(map[string]model.Maintenance)
	if len(maintenanceIds) > 0 {
		maintenanceData, err := f.maintenanceRepo.FindByIds(ctx, maintenanceIds)
		if err != nil {
			return finances, err
		}

		for _, v := range maintenanceData {
			maintenances[v.ID] = v
		}
	}

	userMaintIds := make([]string, 0)
	productsIds := make([]string, 0)
	for _, v := range maintenances {
		if v.ProductID != "" {
			userMaintIds = append(userMaintIds, v.UserID)
			productsIds = append(productsIds, v.ProductID)
		}
	}

	userMaints := make(map[string]model.User)
	if len(userMaintIds) > 0 {
		userMaintData, err := f.userRepo.FindByIds(ctx, userMaintIds)
		if err != nil {
			return finances, err
		}

		for _, v := range userMaintData {
			userMaints[v.ID] = v
		}
	}

	products := make(map[string]model.Product)
	if len(productsIds) > 0 {
		productData, err := f.productRepo.FindByIds(ctx, productsIds)
		if err != nil {
			return finances, err
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
		categorieData, err := f.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return finances, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	users := make(map[string]model.User)
	if len(userIds) > 0 {
		userData, err := f.userRepo.FindByIds(ctx, userIds)
		if err != nil {
			return finances, err
		}

		for _, v := range userData {
			users[v.ID] = v
		}
	}

	financeData := make([]dto.FinanceData, 0)
	for _, v := range results {
		var maintence *dto.MaintenanceData
		if vMaintenance, e := maintenances[v.MaintenanceId]; e {

			var product *dto.ProductData
			if vProduct, e := products[vMaintenance.ProductID]; e {
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

			var userMaint *dto.UserData
			if vUserMaint, e := userMaints[vMaintenance.UserID]; e {
				userMaint = &dto.UserData{
					ID:       vUserMaint.ID,
					Username: vUserMaint.Username,
					Name:     vUserMaint.Name,
					Role:     vUserMaint.Role,
					Status:   vUserMaint.Status,
				}
			}

			maintence = &dto.MaintenanceData{
				ID:          vMaintenance.ID,
				Product:     product,
				Description: vMaintenance.Description,
				Status:      vMaintenance.Status,
				StartDate:   vMaintenance.StartDate.Format("2006-01-02"),
				EndDate:     vMaintenance.EndDate.Format("2006-01-02"),
				Stock:       vMaintenance.Stock,
				User:        userMaint,
			}
		}

		var user *dto.UserData
		if vUser, e := users[v.UserId]; e {
			user = &dto.UserData{
				ID:       vUser.ID,
				Username: vUser.Username,
				Name:     vUser.Name,
				Role:     vUser.Role,
				Status:   vUser.Status,
			}
		}

		financeData = append(financeData, dto.FinanceData{
			ID:          v.ID,
			Maintenance: maintence,
			Description: v.Description,
			Price:       v.Price,
			User:        user,
		})
	}

	return financeData, nil
}

// Search implements [model.FinanceService].
func (f *financeService) Search(ctx context.Context, keyword string, fs model.FinanceSearch) (finances []dto.FinanceData, err error) {
	results, err := f.financeRepo.Search(ctx, keyword, fs)
	if err != nil {
		return finances, err
	}

	maintenanceIds := make([]string, 0)
	userIds := make([]string, 0)
	for _, v := range results {
		maintenanceIds = append(maintenanceIds, v.MaintenanceId)
		userIds = append(userIds, v.UserId)
	}

	maintenances := make(map[string]model.Maintenance)
	if len(maintenanceIds) > 0 {
		maintenanceData, err := f.maintenanceRepo.FindByIds(ctx, maintenanceIds)
		if err != nil {
			return finances, err
		}

		for _, v := range maintenanceData {
			maintenances[v.ID] = v
		}
	}

	userMaintIds := make([]string, 0)
	productsIds := make([]string, 0)
	for _, v := range maintenances {
		if v.ProductID != "" {
			userMaintIds = append(userMaintIds, v.UserID)
			productsIds = append(productsIds, v.ProductID)
		}
	}

	userMaints := make(map[string]model.User)
	if len(userMaintIds) > 0 {
		userMaintData, err := f.userRepo.FindByIds(ctx, userMaintIds)
		if err != nil {
			return finances, err
		}

		for _, v := range userMaintData {
			userMaints[v.ID] = v
		}
	}

	products := make(map[string]model.Product)
	if len(productsIds) > 0 {
		productData, err := f.productRepo.FindByIds(ctx, productsIds)
		if err != nil {
			return finances, err
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
		categorieData, err := f.categorieRepo.FindByIds(ctx, categorieIds)
		if err != nil {
			return finances, err
		}

		for _, v := range categorieData {
			categories[v.ID] = v
		}
	}

	users := make(map[string]model.User)
	if len(userIds) > 0 {
		userData, err := f.userRepo.FindByIds(ctx, userIds)
		if err != nil {
			return finances, err
		}

		for _, v := range userData {
			users[v.ID] = v
		}
	}

	financeData := make([]dto.FinanceData, 0)
	for _, v := range results {
		var maintence *dto.MaintenanceData
		if vMaintenance, e := maintenances[v.MaintenanceId]; e {

			var product *dto.ProductData
			if vProduct, e := products[vMaintenance.ProductID]; e {
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

			var userMaint *dto.UserData
			if vUserMaint, e := userMaints[vMaintenance.UserID]; e {
				userMaint = &dto.UserData{
					ID:       vUserMaint.ID,
					Username: vUserMaint.Username,
					Name:     vUserMaint.Name,
					Role:     vUserMaint.Role,
					Status:   vUserMaint.Status,
				}
			}

			maintence = &dto.MaintenanceData{
				ID:          vMaintenance.ID,
				Product:     product,
				Description: vMaintenance.Description,
				Status:      vMaintenance.Status,
				StartDate:   vMaintenance.StartDate.Format("2006-01-02"),
				EndDate:     vMaintenance.EndDate.Format("2006-01-02"),
				Stock:       vMaintenance.Stock,
				User:        userMaint,
			}
		}

		var user *dto.UserData
		if vUser, e := users[v.UserId]; e {
			user = &dto.UserData{
				ID:       vUser.ID,
				Username: vUser.Username,
				Name:     vUser.Name,
				Role:     vUser.Role,
				Status:   vUser.Status,
			}
		}

		financeData = append(financeData, dto.FinanceData{
			ID:          v.ID,
			Maintenance: maintence,
			Description: v.Description,
			Price:       v.Price,
			User:        user,
		})
	}

	return financeData, nil
}

// Detail implements [model.FinanceService].
func (f *financeService) Detail(ctx context.Context, id string, fs model.FinanceSearch) (finance dto.FinanceData, err error) {
	persisted, err := f.financeRepo.FindById(ctx, id, fs)
	if err != nil {
		return finance, err
	}

	if persisted.ID == "" {
		return finance, errors.New("Data finance tidak ditemukan")
	}

	maintenanceMd, err := f.maintenanceRepo.FindById(
		ctx,
		persisted.MaintenanceId,
		model.MaintenanceSearch{},
	)
	if err != nil {
		return finance, err
	}

	if maintenanceMd.ID == "" {
		return finance, errors.New("Data maintenance tidak ditemukan")
	}

	productMd, err := f.productRepo.FindById(
		ctx,
		maintenanceMd.ProductID,
		model.ProductSearch{},
	)
	if err != nil {
		return finance, err
	}

	if productMd.ID == "" {
		return finance, errors.New("Data product tidak ditemukan")
	}

	userMaintMd, err := f.userRepo.FindById(ctx, maintenanceMd.UserID)
	if err != nil {
		return finance, err
	}

	if userMaintMd.ID == "" {
		return finance, errors.New("Data user maintenance tidak ditemukan")
	}

	var category *dto.CategorieData
	if productMd.CategorieId != "" {
		categories, err := f.categorieRepo.FindByIds(
			ctx,
			[]string{productMd.CategorieId},
		)
		if err != nil {
			return finance, err
		}

		if len(categories) > 0 {
			category = &dto.CategorieData{
				ID:          categories[0].ID,
				Name:        categories[0].Name,
				Description: categories[0].Description,
			}
		}
	}

	product := &dto.ProductData{
		ID:        productMd.ID,
		Name:      productMd.Name,
		Categorie: category,
		Merek:     productMd.Merek,
	}

	userMaint := &dto.UserData{
		ID:       userMaintMd.ID,
		Username: userMaintMd.Username,
		Name:     userMaintMd.Name,
		Role:     userMaintMd.Role,
		Status:   userMaintMd.Status,
	}

	maintenance := &dto.MaintenanceData{
		ID:          maintenanceMd.ID,
		Product:     product,
		Description: maintenanceMd.Description,
		Status:      maintenanceMd.Status,
		StartDate:   maintenanceMd.StartDate.Format("2006-01-02"),
		EndDate:     maintenanceMd.EndDate.Format("2006-01-02"),
		Stock:       maintenanceMd.Stock,
		User:        userMaint,
	}

	userMd, err := f.userRepo.FindById(ctx, persisted.UserId)
	if err != nil {
		return finance, err
	}

	if userMd.ID == "" {
		return finance, errors.New("Data user finance tidak ditemukan")
	}

	user := &dto.UserData{
		ID:       userMd.ID,
		Username: userMd.Username,
		Name:     userMd.Name,
		Role:     userMd.Role,
		Status:   userMd.Status,
	}

	return dto.FinanceData{
		ID:          persisted.ID,
		Maintenance: maintenance,
		Description: persisted.Description,
		Price:       persisted.Price,
		User:        user,
	}, nil
}

// Create implements [model.FinanceService].
func (f *financeService) Create(ctx context.Context, userId string, maintenanceId string, req dto.FinanceCreated) error {
	persisted, err := f.maintenanceRepo.FindById(ctx, maintenanceId, model.MaintenanceSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data maintenace tidak ditemukan")
	}

	finance := model.Finance{
		ID:            uuid.NewString(),
		MaintenanceId: persisted.ID,
		Description:   req.Description,
		Price:         req.Price,
		UserId:        userId,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	return f.financeRepo.Save(ctx, &finance)
}

// Update implements [model.FinanceService].
func (f *financeService) Update(ctx context.Context, req dto.FinanceUpdated) error {
	persisted, err := f.financeRepo.FindById(ctx, req.ID, model.FinanceSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	persisted.Description = req.Description
	persisted.Price = req.Price
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	return f.financeRepo.Update(ctx, &persisted)
}

// Delete implements [model.FinanceService].
func (f *financeService) Delete(ctx context.Context, id string) error {
	persisted, err := f.financeRepo.FindById(ctx, id, model.FinanceSearch{})
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	return f.financeRepo.Delete(ctx, id)
}
