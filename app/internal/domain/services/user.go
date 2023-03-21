package services

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"jungle-test/app/internal/domain/entity"
	"jungle-test/app/pkg/apperrors"
)

type UserStorage interface {
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type UserService struct {
	userStorage UserStorage
}

func NewUserService(userStorage UserStorage) *UserService {
	return &UserService{userStorage: userStorage}
}

func (s UserService) Login(ctx context.Context, username string, password string) (*entity.User, error) {

	user, err := s.userStorage.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)) != nil {
		return nil, apperrors.ErrWrongPassword
	}

	user.PasswordHash = nil
	return user, nil
}
