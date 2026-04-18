package model

import (
	"context"
	"database/sql"
	"inventory-app/internal/dto"
)

type User struct {
	ID        string       `db:"id"`
	Username  string       `db:"username"`
	Name      string       `db:"name"`
	Role      string       `db:"role"`
	Password  string       `db:"password"`
	Status    string       `db:"status"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type UserRepository interface {
	FindAll(ctx context.Context) (users []User, err error)
	FindById(ctx context.Context, id string) (user User, err error)
	FindByIds(ctx context.Context, ids []string) (users []User, err error)
	FindByUsername(ctx context.Context, username string) (user User, err error)
	Search(ctx context.Context, keyword string) (users []User, err error)
	Save(ctx context.Context, us *User) error
	UpdateData(ctx context.Context, us *User) error
	UpdatePassword(ctx context.Context, us *User) error
	Delete(ctx context.Context, id string) error
}

type UserService interface {
	Index(ctx context.Context) (users []dto.UserData, err error)
	Search(ctx context.Context, keyword string) (users []dto.UserData, err error)
	Detail(ctx context.Context, id string) (user dto.UserData, err error)
	Create(ctx context.Context, req dto.UserCreated) error
	UpdateData(ctx context.Context, req dto.UserUpdateData) error
	UpdatePassword(ctx context.Context, req dto.UserUpdatePassword) error
	Delete(ctx context.Context, id string) error
}

type UserAuthRepository interface {
	FindByUsername(ctx context.Context, username string) (user User, err error)
	HandleLoginResult(ctx context.Context, username string, success bool) error
}
