package service

import (
	"context"
	"database/sql"
	"errors"
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo model.UserRepository
}

func NewUser(userRepo model.UserRepository) model.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Index implements [model.UserService].
func (u *userService) Index(ctx context.Context) (users []dto.UserData, err error) {
	results, err := u.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		users = append(users, dto.UserData{
			ID:       v.ID,
			Username: v.Username,
			Name:     v.Name,
			Role:     v.Role,
			Status:   v.Status,
		})
	}

	return users, nil
}

// Search implements [model.UserService].
func (u *userService) Search(ctx context.Context, keyword string) (users []dto.UserData, err error) {
	results, err := u.userRepo.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		users = append(users, dto.UserData{
			ID:       v.ID,
			Username: v.Username,
			Name:     v.Name,
			Role:     v.Role,
			Status:   v.Status,
		})
	}

	return users, nil
}

// Detail implements [model.UserService].
func (u *userService) Detail(ctx context.Context, id string) (user dto.UserData, err error) {
	persisted, err := u.userRepo.FindById(ctx, id)
	if err != nil {
		return user, err
	}

	if persisted.ID == "" {
		return user, errors.New("Data tidak ditemukan")
	}

	return dto.UserData{
		ID:       persisted.ID,
		Username: persisted.Username,
		Name:     persisted.Name,
		Role:     persisted.Role,
		Status:   persisted.Status,
	}, nil
}

// Create implements [model.UserService].
func (u *userService) Create(ctx context.Context, req dto.UserCreated) error {
	exist, err := u.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	if exist.ID != "" {
		return errors.New("Username telah digunakan")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Name:     req.Name,
		Role:     req.Role,
		Password: string(hashPassword),
		Status:   model.UserStatusAktif,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	return u.userRepo.Save(ctx, &user)
}

// UpdateData implements [model.UserService].
func (u *userService) UpdateData(ctx context.Context, req dto.UserUpdateData) error {
	persisted, err := u.userRepo.FindById(ctx, req.ID)
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	persisted.Name = req.Name
	persisted.Role = req.Role
	persisted.Status = req.Status
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	return u.userRepo.UpdateData(ctx, &persisted)
}

// UpdatePassword implements [model.UserService].
func (u *userService) UpdatePassword(ctx context.Context, req dto.UserUpdatePassword) error {
	persisted, err := u.userRepo.FindById(ctx, req.ID)
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	persisted.Password = string(hashPassword)
	persisted.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	return u.userRepo.UpdatePassword(ctx, &persisted)
}

// Delete implements [model.UserService].
func (u *userService) Delete(ctx context.Context, id string) error {
	persisted, err := u.userRepo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if persisted.ID == "" {
		return errors.New("Data tidak ditemukan")
	}

	return u.userRepo.Delete(ctx, id)
}
