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

type warehouseRepository struct {
	db    *goqu.Database
	cache *cache.Cache
}

func NewWarehouse(con *sql.DB, rdb *redis.Client, env string) model.WarehouseRepository {
	return &warehouseRepository{
		db:    goqu.New("mysql", con),
		cache: cache.New(rdb, env+":warehouse", 5*time.Minute, true),
	}
}

// FindAll implements [model.WarehouseRepository].
func (w *warehouseRepository) FindAll(ctx context.Context, ws model.WarehouseSearch) (warehouses []model.Warehouse, err error) {
	cacheKey := fmt.Sprintf(
		"warehouses:product=%s",
		ws.ProductId,
	)

	err = w.cache.Get(ctx, cacheKey, &warehouses)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := w.db.From("warehouses").Where(goqu.C("deleted_at").IsNull())
	if ws.ProductId != "" {
		dataset = dataset.Where(goqu.C("product_id").Eq(ws.ProductId))
	}

	err = dataset.ScanStructsContext(ctx, &warehouses)
	if err != nil {
		return warehouses, err
	}

	if len(warehouses) > 0 {
		w.cache.Set(ctx, cacheKey, warehouses, 3*time.Minute)
	}

	return
}

// FindById implements [model.WarehouseRepository].
func (w *warehouseRepository) FindById(ctx context.Context, id string, ws model.WarehouseSearch) (warehouse model.Warehouse, err error) {
	cacheKey := fmt.Sprintf(
		"warehouses:%s:product=%s",
		id,
		ws.ProductId,
	)

	err = w.cache.Get(ctx, cacheKey, &warehouse)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := w.db.From("warehouses").Where(goqu.C("deleted_at").IsNull())
	if ws.ProductId != "" {
		dataset = dataset.Where(goqu.C("product_id").Eq(ws.ProductId))
	}

	_, err = dataset.ScanStructContext(ctx, &warehouse)
	if err != nil {
		return warehouse, err
	}

	w.cache.Set(ctx, cacheKey, warehouse, 3*time.Minute)

	return
}

// FindByIds implements [model.WarehouseRepository].
func (w *warehouseRepository) FindByIds(ctx context.Context, ids []string) (warehouses []model.Warehouse, err error) {
	if len(ids) == 0 {
		return warehouses, nil
	}

	dataset := w.db.From("warehouses").Where(goqu.C("id").In(ids))
	err = dataset.ScanStructsContext(ctx, &warehouses)
	if err != nil {
		return warehouses, err
	}

	return
}

// FindByProductId implements [model.WarehouseRepository].
func (w *warehouseRepository) FindByProductId(ctx context.Context, productId string) (warehouse model.Warehouse, err error) {
	dataset := w.db.From("warehouses").Where(goqu.C("product_id").Eq(productId))

	_, err = dataset.ScanStructContext(ctx, &warehouse)
	if err != nil {
		return warehouse, err
	}

	return
}

// Search implements [model.WarehouseRepository].
func (w *warehouseRepository) Search(ctx context.Context, keyword string, ws model.WarehouseSearch) (warehouses []model.Warehouse, err error) {
	cacheKey := fmt.Sprintf(
		"warehouses:%s:product=%s",
		keyword,
		ws.ProductId,
	)

	err = w.cache.Get(ctx, cacheKey, &warehouses)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := w.db.Select("warehouses.*").From("warehouses").Join(
		goqu.T("products"),
		goqu.On(goqu.I("warehouses.product_id").Eq(goqu.I("products.id"))),
	).Where(
		goqu.I("warehouses.deleted_at").IsNull(),
		goqu.Or(
			goqu.I("warehouses.status").Like("%"+keyword+"%"),
			goqu.I("products.name").Like("%"+keyword+"%"),
			goqu.I("products.merek").Like("%"+keyword+"%"),
		),
	)
	if ws.ProductId != "" {
		dataset = dataset.Where(goqu.I("warehouses.product_id").Eq(ws.ProductId))
	}

	err = dataset.ScanStructsContext(ctx, &warehouses)
	if err != nil {
		return warehouses, err
	}

	if len(warehouses) > 0 {
		w.cache.Set(ctx, cacheKey, warehouses, 3*time.Minute)
	}

	return
}

// Save implements [model.WarehouseRepository].
func (w *warehouseRepository) Save(ctx context.Context, wh *model.Warehouse) error {
	_, err := w.db.Insert("warehouses").Rows(wh).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	w.cache.DeleteByPrefix(ctx, "")

	return err
}

// Update implements [model.WarehouseRepository].
func (w *warehouseRepository) Update(ctx context.Context, wh *model.Warehouse) error {
	_, err := w.db.Update("warehouses").Where(goqu.C("id").Eq(wh.ID)).Set(wh).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	w.cache.DeleteByPrefix(ctx, "")

	return err
}

// Delete implements [model.WarehouseRepository].
func (w *warehouseRepository) Delete(ctx context.Context, id string) error {
	_, err := w.db.Update("warehouses").Where(goqu.C("id").Eq(id)).Set(
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

	w.cache.DeleteByPrefix(ctx, "")

	return err
}
