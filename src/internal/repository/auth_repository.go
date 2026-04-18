package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory-app/internal/cache"
	"inventory-app/internal/model"
	"inventory-app/internal/util"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/redis/go-redis/v9"
)

type authRepository struct {
	db    *goqu.Database
	cache *cache.Cache
}

func NewAuth(con *sql.DB, rdb *redis.Client, env string) model.UserAuthRepository {
	return &authRepository{
		db:    goqu.New("mysql", con),
		cache: cache.New(rdb, env+":login", 5*time.Minute, false),
	}
}

// FindByUsername implements [model.UserAuthRepository].
func (a *authRepository) FindByUsername(ctx context.Context, username string) (user model.User, err error) {
	var blocked bool

	err = a.cache.Get(ctx, util.BlockKey(username), &blocked)

	if err == nil {
		return user, errors.New("Terlalu banyak percobaan login")
	}

	if err != cache.ErrCacheMiss {
		return user, err
	}

	dataset := a.db.From("users").Where(goqu.I("username").Eq(username))
	_, err = dataset.ScanStructContext(ctx, &user)

	return
}

// HandleLoginResult implements [model.UserAuthRepository].
func (a *authRepository) HandleLoginResult(ctx context.Context, username string, success bool) error {
	attemptKey := util.AttempKey(username)
	blockKey := util.BlockKey(username)

	if success {
		if err := a.cache.Delete(ctx, attemptKey); err != nil {
			return err
		}

		if err := a.cache.Delete(ctx, blockKey); err != nil {
			return err
		}

		return nil
	}

	var attempts int
	err := a.cache.Get(ctx, attemptKey, &attempts)

	switch err {
	case cache.ErrCacheMiss:
		attempts = 1
	case nil:
		attempts++
	default:
		return err
	}

	_ = a.cache.Set(ctx, attemptKey, attempts, 5*time.Minute)

	if attempts >= 3 {
		_ = a.cache.Set(ctx, blockKey, true, 5*time.Minute)
	}

	return nil
}
