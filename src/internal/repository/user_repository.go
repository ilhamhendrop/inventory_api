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

type userRepository struct {
	db    *goqu.Database
	cache *cache.Cache
}

func NewUser(con *sql.DB, rdb *redis.Client, env string) model.UserRepository {
	return &userRepository{
		db:    goqu.New("mysql", con),
		cache: cache.New(rdb, env+":user", 5*time.Minute, true),
	}
}

// FindAll implements [model.UserRepository].
func (u *userRepository) FindAll(ctx context.Context) (users []model.User, err error) {
	cacheKey := "users"

	err = u.cache.Get(ctx, cacheKey, &users)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := u.db.From("users").Where(goqu.C("deleted_at").IsNull())
	err = dataset.ScanStructsContext(ctx, &users)
	if err != nil {
		return users, err
	}

	if len(users) > 0 {
		u.cache.Set(ctx, cacheKey, users, 3*time.Minute)
	}

	return
}

// FindById implements [model.UserRepository].
func (u *userRepository) FindById(ctx context.Context, id string) (user model.User, err error) {
	cacheKey := fmt.Sprintf(
		"users:%s",
		id,
	)

	err = u.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := u.db.From("users").Where(goqu.C("id").Eq(id))
	_, err = dataset.ScanStructContext(ctx, &user)
	if err != nil {
		return user, err
	}

	_ = u.cache.Set(ctx, cacheKey, user, 3*time.Minute)

	return
}

// FindByIds implements [model.UserRepository].
func (u *userRepository) FindByIds(ctx context.Context, ids []string) (users []model.User, err error) {
	if len(ids) == 0 {
		return users, nil
	}

	dataset := u.db.From("users").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").In(ids))
	err = dataset.ScanStructsContext(ctx, &users)
	if err != nil {
		return users, err
	}

	return
}

// FindByUsername implements [model.UserRepository].
func (u *userRepository) FindByUsername(ctx context.Context, username string) (user model.User, err error) {
	dataset := u.db.From("users").Where(goqu.C("deleted_at").IsNull(), goqu.C("username").Eq(username)).Limit(1)
	_, err = dataset.ScanStructContext(ctx, &user)
	if err != nil {
		return user, err
	}

	return
}

// Search implements [model.UserRepository].
func (u *userRepository) Search(ctx context.Context, keyword string) (users []model.User, err error) {
	cacheKey := fmt.Sprintf(
		"users:search=%s",
		keyword,
	)

	err = u.cache.Get(ctx, cacheKey, &users)
	if err == nil {
		return
	}

	if err != cache.ErrCacheMiss {
		return
	}

	dataset := u.db.From("users").Where(goqu.C("deleted_at").IsNull(), goqu.Or(
		goqu.C("username").Like("%"+keyword+"%"),
		goqu.C("name").Like("%"+keyword+"%"),
		goqu.C("status").Like("%"+keyword+"%"),
		goqu.C("role").Like("%"+keyword+"%"),
	))

	err = dataset.ScanStructsContext(ctx, &users)
	if err != nil {
		return users, err
	}

	if len(users) > 0 {
		u.cache.Set(ctx, cacheKey, users, 3*time.Minute)
	}

	return
}

// Save implements [model.UserRepository].
func (u *userRepository) Save(ctx context.Context, us *model.User) error {
	_, err := u.db.Insert("users").Rows(us).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	u.cache.DeleteByPrefix(ctx, "")

	return nil
}

// UpdateData implements [model.UserRepository].
func (u *userRepository) UpdateData(ctx context.Context, us *model.User) error {
	_, err := u.db.Update("users").Where(goqu.C("id").Eq(us.ID)).Set(us).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	u.cache.DeleteByPrefix(ctx, "")

	return nil
}

// UpdatePassword implements [model.UserRepository].
func (u *userRepository) UpdatePassword(ctx context.Context, us *model.User) error {
	_, err := u.db.Update("users").Where(goqu.C("id").Eq(us.ID)).Set(us).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	u.cache.DeleteByPrefix(ctx, "")

	return nil
}

// Delete implements [model.UserRepository].
func (u *userRepository) Delete(ctx context.Context, id string) error {
	_, err := u.db.Update("users").Where(goqu.C("id").Eq(id)).Set(goqu.Record{
		"deleted_at": sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}).Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	u.cache.DeleteByPrefix(ctx, "")

	return nil
}
