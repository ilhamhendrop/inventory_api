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

type financeRepository struct {
	db    *goqu.Database
	cache *cache.Cache
}

func NewFinance(con *sql.DB, rdb *redis.Client, env string) model.FinanceRepository {
	return &financeRepository{
		db:    goqu.New("mysql", con),
		cache: cache.New(rdb, env+":finance", 5*time.Minute, true),
	}
}

// FindAll implements [model.FinanceRepository].
func (f *financeRepository) FindAll(ctx context.Context, fs model.FinanceSearch) (finances []model.Finance, err error) {
	cacheKey := fmt.Sprintf(
		"finances:maintenance=%s:user=%s",
		fs.MaintenanceId,
		fs.UserId,
	)

	err = f.cache.Get(ctx, cacheKey, &finances)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := f.db.From("finances").Where(goqu.C("deleted_at").IsNull())
	if fs.MaintenanceId != "" {
		dataset = dataset.Where(goqu.C("maintenance_id").Eq(fs.MaintenanceId))
	}

	if fs.UserId != "" {
		dataset = dataset.Where(goqu.C("user_id").Eq(fs.UserId))
	}

	err = dataset.ScanStructsContext(ctx, &finances)
	if err != nil {
		return finances, err
	}

	if len(finances) > 0 {
		f.cache.Set(ctx, cacheKey, finances, 3*time.Minute)
	}

	return
}

// FindById implements [model.FinanceRepository].
func (f *financeRepository) FindById(ctx context.Context, id string, fs model.FinanceSearch) (finance model.Finance, err error) {
	cacheKey := fmt.Sprintf(
		"finance:%s:maintenance=%s:user=%s",
		id,
		fs.MaintenanceId,
		fs.UserId,
	)

	err = f.cache.Get(ctx, cacheKey, &finance)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := f.db.From("finances").Where(goqu.C("deleted_at").IsNull())
	if fs.MaintenanceId != "" {
		dataset = dataset.Where(goqu.C("maintenance_id").Eq(fs.MaintenanceId))
	}

	if fs.UserId != "" {
		dataset = dataset.Where(goqu.C("user_id").Eq(fs.UserId))
	}

	_, err = dataset.ScanStructContext(ctx, &finance)

	f.cache.Set(ctx, cacheKey, finance, 3*time.Minute)

	return
}

// Search implements [model.FinanceRepository].
func (f *financeRepository) Search(ctx context.Context, keyword string, fs model.FinanceSearch) (finances []model.Finance, err error) {
	cacheKey := fmt.Sprintf(
		"finances:maintenance=%s:user=%s:search=%s",
		fs.MaintenanceId,
		fs.UserId,
		keyword,
	)

	err = f.cache.Get(ctx, cacheKey, &finances)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := f.db.Select("finances.*").From("finances").Join(
		goqu.T("maintenances"),
		goqu.On(goqu.I("finances.maintenance_id").Eq(goqu.I("maintenances.id"))),
	).Join(
		goqu.T("users"),
		goqu.On(goqu.I("finances.user_id").Eq(goqu.I("users.id"))),
	).Where(
		goqu.I("maintenances.deleted_at").IsNull(),
		goqu.Or(
			goqu.I("maintenances.start_date").Like("%"+keyword+"%"),
			goqu.I("maintenances.end_date").Like("%"+keyword+"%"),
			goqu.I("users.name").Like("%"+keyword+"%"),
		),
	)

	if fs.MaintenanceId != "" {
		dataset = dataset.Where(goqu.I("finances.maintenance_id").Eq(fs.MaintenanceId))
	}

	if fs.UserId != "" {
		dataset = dataset.Where(goqu.I("finances.user_id").Eq(fs.UserId))
	}

	err = dataset.ScanStructsContext(ctx, &finances)
	if err != nil {
		return finances, err
	}

	if len(finances) > 0 {
		f.cache.Set(ctx, cacheKey, finances, 3*time.Minute)
	}

	return
}

// Save implements [model.FinanceRepository].
func (f *financeRepository) Save(ctx context.Context, fc *model.Finance) error {
	_, err := f.db.Insert("finances").Rows(fc).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	f.cache.DeleteByPrefix(ctx, "")

	return nil
}

// Update implements [model.FinanceRepository].
func (f *financeRepository) Update(ctx context.Context, fc *model.Finance) error {
	_, err := f.db.Update("finances").Where(
		goqu.C("id").Eq(fc.ID)).Set(fc).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	f.cache.DeleteByPrefix(ctx, "")

	return nil
}

// Delete implements [model.FinanceRepository].
func (f *financeRepository) Delete(ctx context.Context, id string) error {
	_, err := f.db.Update("finances").Where(goqu.C("id").Eq(id)).Set(
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

	f.cache.DeleteByPrefix(ctx, "")

	return nil
}
