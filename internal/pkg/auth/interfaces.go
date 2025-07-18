package auth

import (
	"context"
	"errors"

	"github.com/K1tten2005/go_vk_intern/internal/models"
)

var (
	ErrInvalidPassword    = errors.New("incorrect password format")
	ErrInvalidLogin       = errors.New("incorrect login format")
	ErrInvalidCredentials = errors.New("wrong login or password")
	ErrCreatingUser       = errors.New("user creation error")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this login already exists")
	ErrGeneratingToken    = errors.New("token generation error")
)

type AuthUsecase interface {
	SignIn(ctx context.Context, data models.UserReq) (models.UserResp, string, error)
	SignUp(ctx context.Context, data models.UserReq) (models.UserResp, string, error)
}

type AuthRepo interface {
	InsertUser(ctx context.Context, user models.User) error
	SelectUserByLogin(ctx context.Context, login string) (models.User, error)
}
