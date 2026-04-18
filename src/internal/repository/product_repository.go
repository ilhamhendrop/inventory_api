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

type productRepository struct {
	db    *goqu.Database
	cache *cache.Cache
}

func NewProduct(con *sql.DB, rdb *redis.Client, env string) model.ProductRepository {
	return &productRepository{
		db:    goqu.New("mysql", con),
		cache: cache.New(rdb, env+":product", 5*time.Minute, true),
	}
}

// FindAll implements [model.ProductRepository].
func (p *productRepository) FindAll(ctx context.Context, ps model.ProductSearch) (products []model.Product, err error) {
	cacheKey := fmt.Sprintf(
		"products:categorie=%s",
		ps.CategorieId,
	)

	err = p.cache.Get(ctx, cacheKey, &products)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := p.db.From("products").Where(goqu.C("deleted_at").IsNull())
	if ps.CategorieId != "" {
		dataset = dataset.Where(goqu.C("categorie_id").Eq(ps.CategorieId))
	}

	err = dataset.ScanStructsContext(ctx, &products)
	if err != nil {
		return products, err
	}

	if len(products) > 0 {
		p.cache.DeleteByPrefix(ctx, "")
	}

	return
}

// FindById implements [model.ProductRepository].
func (p *productRepository) FindById(ctx context.Context, id string, ps model.ProductSearch) (product model.Product, err error) {
	cacheKey := fmt.Sprintf(
		"products:%s:categorie=%s",
		id,
		ps.CategorieId,
	)

	err = p.cache.Get(ctx, cacheKey, &product)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := p.db.From("products").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").Eq(id))
	if ps.CategorieId != "" {
		dataset = dataset.Where(goqu.C("categorie_id").Eq(ps.CategorieId))
	}

	_, err = dataset.ScanStructContext(ctx, &product)
	if err != nil {
		return product, err
	}

	p.cache.Set(ctx, cacheKey, product, 3*time.Minute)

	return
}

// FindByIds implements [model.ProductRepository].
func (p *productRepository) FindByIds(ctx context.Context, ids []string) (products []model.Product, err error) {
	if len(ids) == 0 {
		return products, nil
	}

	dataset := p.db.From("products").Where(goqu.C("id").In(ids))
	err = dataset.ScanStructsContext(ctx, &products)
	if err != nil {
		return products, err
	}

	return
}

// Search implements [model.ProductRepository].
func (p *productRepository) Search(ctx context.Context, keyword string, ps model.ProductSearch) (products []model.Product, err error) {
	cacheKey := fmt.Sprintf(
		"products:search=%s:categorie=%s",
		keyword,
		ps.CategorieId,
	)

	err = p.cache.Get(ctx, cacheKey, &products)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := p.db.Select("products.*").From("products").Join(
		goqu.I("categories"),
		goqu.On(goqu.I("products.categorie_id").Eq(goqu.I("categories.id"))),
	).Where(
		goqu.I("products.deleted_at").IsNull(),
		goqu.Or(
			goqu.I("products.name").Like("%"+keyword+"%"),
			goqu.I("products.merek").Like("%"+keyword+"%"),
			goqu.I("categories.name").Like("%"+keyword+"%"),
		),
	)
	if ps.CategorieId != "" {
		dataset = dataset.Where(goqu.I("products.categorie_id").Eq(ps.CategorieId))
	}

	err = dataset.ScanStructsContext(ctx, &products)
	if err != nil {
		return products, err
	}

	if len(products) > 0 {
		p.cache.Set(ctx, cacheKey, products, 3*time.Minute)
	}

	return
}

// Save implements [model.ProductRepository].
func (p *productRepository) Save(ctx context.Context, pr *model.Product) error {
	_, err := p.db.Insert("products").Rows(pr).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	p.cache.DeleteByPrefix(ctx, "")

	return err
}

// Update implements [model.ProductRepository].
func (p *productRepository) Update(ctx context.Context, pr *model.Product) error {
	_, err := p.db.Update("products").Where(
		goqu.C("id").Eq(pr.ID),
	).Set(pr).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	p.cache.DeleteByPrefix(ctx, "")

	return err
}

// Delete implements [model.ProductRepository].
func (p *productRepository) Delete(ctx context.Context, id string) error {
	_, err := p.db.Update("products").Where(
		goqu.C("id").Eq(id),
	).Set(
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

	p.cache.DeleteByPrefix(ctx, "")

	return err
}
