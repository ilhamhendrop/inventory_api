package repository

import (
	"context"
	"database/sql"
	"fmt"
	"inventory-app/internal/cache"
	"inventory-app/internal/model"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/redis/go-redis/v9"
)

type categorieRepository struct {
	db    *goqu.Database
	cache *cache.Cache
}

func NewCategorie(con *sql.DB, rdb *redis.Client, env string) model.CategorieRepository {
	return &categorieRepository{
		db:    goqu.New("mysql", con),
		cache: cache.New(rdb, env+":categorie", 5*time.Minute, true),
	}
}

// FindAll implements [model.CategorieRepository].
func (c *categorieRepository) FindAll(ctx context.Context) (categries []model.Categorie, err error) {
	cacheKey := "categories"

	err = c.cache.Get(ctx, cacheKey, &categries)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := c.db.From("categories").Where(goqu.C("deleted_at").IsNull())
	err = dataset.ScanStructsContext(ctx, &categries)
	if err != nil {
		return categries, err
	}

	if len(categries) > 0 {
		c.cache.Set(ctx, cacheKey, categries, 3*time.Minute)
	}

	return
}

// FindById implements [model.CategorieRepository].
func (c *categorieRepository) FindById(ctx context.Context, id string) (categorie model.Categorie, err error) {
	cacheKey := fmt.Sprintf(
		"categories:%s",
		id,
	)

	err = c.cache.Get(ctx, cacheKey, &categorie)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := c.db.From("categories").Where(goqu.C("id").Eq(id))
	_, err = dataset.ScanStructContext(ctx, &categorie)
	if err != nil {
		return categorie, err
	}

	_ = c.cache.Set(ctx, cacheKey, categorie, 3*time.Minute)

	return
}

// FindByIds implements [model.CategorieRepository].
func (c *categorieRepository) FindByIds(ctx context.Context, ids []string) (categories []model.Categorie, err error) {
	if len(ids) == 0 {
		return categories, nil
	}

	dataset := c.db.From("categories").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").In(ids))
	err = dataset.ScanStructsContext(ctx, &categories)
	if err != nil {
		return categories, err
	}

	return
}

// Search implements [model.CategorieRepository].
func (c *categorieRepository) Search(ctx context.Context, keyword string) (categories []model.Categorie, err error) {
	cacheKey := fmt.Sprintf(
		"categories:search=%s",
		keyword,
	)

	err = c.cache.Get(ctx, cacheKey, &categories)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := c.db.From("categories").Where(goqu.C("deleted_at").IsNull(), goqu.Or(
		goqu.C("name").Like("%"+keyword+"%"),
	))
	err = dataset.ScanStructsContext(ctx, &categories)
	if err != nil {
		return categories, err
	}

	if len(categories) > 0 {
		c.cache.Set(ctx, cacheKey, categories, 3*time.Minute)
	}

	return
}

// Save implements [model.CategorieRepository].
func (c *categorieRepository) Save(ctx context.Context, ct *model.Categorie) error {
	_, err := c.db.Insert("categories").Rows(ct).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	c.cache.DeleteByPrefix(ctx, "")

	return nil
}

// Update implements [model.CategorieRepository].
func (c *categorieRepository) Update(ctx context.Context, ct *model.Categorie) error {
	_, err := c.db.Update("categories").Where(goqu.C("id").Eq(ct.ID)).Set(ct).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	c.cache.DeleteByPrefix(ctx, "")

	return nil
}

// Delete implements [model.CategorieRepository].
func (c *categorieRepository) Delete(ctx context.Context, id string) error {
	_, err := c.db.Update("categories").Where(goqu.C("id").Eq(id)).Set(goqu.Record{
		"deleted_at": sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}).Executor().ExecContext(ctx)

	if err != nil {
		return err
	}

	c.cache.DeleteByPrefix(ctx, "")

	return nil
}
