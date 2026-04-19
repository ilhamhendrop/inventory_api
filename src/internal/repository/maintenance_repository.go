package repository

import (
	"context"
	"database/sql"
	"fmt"
	"inventory-app/internal/cache"
	"inventory-app/internal/model"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/redis/go-redis/v9"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

type maintenanceRepository struct {
	db    *goqu.Database
	cache *cache.Cache
}

func NewMaintenance(con *sql.DB, rdb *redis.Client, env string) model.MaintenanceRepository {
	return &maintenanceRepository{
		db:    goqu.New("mysql", con),
		cache: cache.New(rdb, env+":maintenance", 5*time.Minute, true),
	}
}

// FindAll implements [model.MaintenanceRepository].
func (m *maintenanceRepository) FindAll(ctx context.Context, ms model.MaintenanceSearch) (maintenances []model.Maintenance, err error) {
	cacheKey := fmt.Sprintf(
		"maintenances:product=%s:user=%s",
		ms.ProductID,
		ms.UserID,
	)

	err = m.cache.Get(ctx, cacheKey, &maintenances)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := m.db.From("maintenances").Where(goqu.C("deleted_at").IsNull())
	if ms.ProductID != "" {
		dataset = dataset.Where(goqu.C("product_id").Eq(ms.ProductID))
	}

	if ms.UserID != "" {
		dataset = dataset.Where(goqu.C("user_id").Eq(ms.UserID))
	}

	err = dataset.ScanStructsContext(ctx, &maintenances)
	if err != nil {
		return maintenances, err
	}

	if len(maintenances) > 0 {
		m.cache.Set(ctx, cacheKey, maintenances, 3*time.Minute)
	}

	return
}

// FindById implements [model.MaintenanceRepository].
func (m *maintenanceRepository) FindById(ctx context.Context, id string, ms model.MaintenanceSearch) (maintenance model.Maintenance, err error) {
	cacheKey := fmt.Sprintf(
		"maintenances:%s:product=%s:user=%s",
		id,
		ms.ProductID,
		ms.UserID,
	)

	err = m.cache.Get(ctx, cacheKey, &maintenance)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := m.db.From("maintenances").Where(goqu.C("deleted_at").IsNull())

	if ms.ProductID != "" {
		dataset = dataset.Where(goqu.C("product_id").Eq(ms.ProductID))
	}

	if ms.UserID != "" {
		dataset = dataset.Where(goqu.C("user_id").Eq(ms.UserID))
	}

	_, err = dataset.ScanStructContext(ctx, &maintenance)
	if err != nil {
		return maintenance, err
	}

	m.cache.Set(ctx, cacheKey, maintenance, 3*time.Minute)

	return
}

// FindByIds implements [model.MaintenanceRepository].
func (m *maintenanceRepository) FindByIds(ctx context.Context, ids []string) (maintenances []model.Maintenance, err error) {
	if len(ids) == 0 {
		return maintenances, nil
	}

	dataset := m.db.From("maintenances").Where(goqu.C("id").In(ids))
	err = dataset.ScanStructsContext(ctx, &maintenances)
	if err != nil {
		return maintenances, err
	}

	return
}

// Search implements [model.MaintenanceRepository].
func (m *maintenanceRepository) Search(ctx context.Context, keyword string, ms model.MaintenanceSearch) (maintenance []model.Maintenance, err error) {
	cacheKey := fmt.Sprintf(
		"maintenances:%s:product=%s",
		keyword,
		ms.ProductID,
	)

	err = m.cache.Get(ctx, cacheKey, &maintenance)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := m.db.Select("maintenances.*").From("maintenances").Join(
		goqu.T("products"),
		goqu.On(goqu.I("maintenances.product_id").Eq(goqu.I("products.id"))),
	).Join(
		goqu.T("users"),
		goqu.On(goqu.I("maintenances.user_id").Eq(goqu.I("users.id"))),
	).Where(
		goqu.I("maintenances.deleted_at").IsNull(),
		goqu.Or(
			goqu.I("maintenances.status").Like("%"+keyword+"%"),
			goqu.I("products.name").Like("%"+keyword+"%"),
			goqu.I("products.merek").Like("%"+keyword+"%"),
			goqu.I("users.name").Like("%"+keyword+"%"),
		),
	)

	if ms.ProductID != "" {
		dataset = dataset.Where(goqu.I("maintenances.products").Eq(ms.ProductID))
	}

	if ms.UserID != "" {
		dataset = dataset.Where(goqu.C("user_id").Eq(ms.UserID))
	}

	err = dataset.ScanStructsContext(ctx, &maintenance)
	if err != nil {
		return maintenance, err
	}

	if len(maintenance) > 0 {
		m.cache.Set(ctx, cacheKey, maintenance, 3*time.Minute)
	}

	return
}

// Save implements [model.MaintenanceRepository].
func (m *maintenanceRepository) Save(ctx context.Context, mt *model.Maintenance) error {
	_, err := m.db.Insert("maintenances").Rows(mt).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	m.cache.DeleteByPrefix(ctx, "")
	return nil
}

// Update implements [model.MaintenanceRepository].
func (m *maintenanceRepository) Update(ctx context.Context, mt *model.Maintenance) error {
	_, err := m.db.Update("maintenances").Where(goqu.C("id").Eq(mt.ID)).Set(mt).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	m.cache.DeleteByPrefix(ctx, "")

	return nil
}

// Delete implements [model.MaintenanceRepository].
func (m *maintenanceRepository) Delete(ctx context.Context, id string) error {
	_, err := m.db.Update("maintenances").Where(goqu.C("id").Eq(id)).Set(
		goqu.Record{
			"deleted_at": sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
		},
	).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	m.cache.DeleteByPrefix(ctx, "")

	return nil
}
