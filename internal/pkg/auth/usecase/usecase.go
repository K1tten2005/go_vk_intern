package usecase

import (
	"context"
	"crypto/rand"
	"log/slog"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/auth"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/jwtUtils"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/validation"
	"github.com/satori/uuid"
)

type AuthUsecase struct {
	repo auth.AuthRepo
}

func CreateAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (uc *AuthUsecase) SignIn(ctx context.Context, data models.UserReq) (models.UserResp, string, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	user, err := uc.repo.SelectUserByLogin(ctx, data.Login)
	if err != nil {
		loggerVar.Error(auth.ErrUserNotFound.Error())
		return models.UserResp{}, "", auth.ErrUserNotFound
	}

	if !validation.CheckPassword(user.PasswordHash, data.Password) {
		loggerVar.Error(auth.ErrInvalidCredentials.Error())
		return models.UserResp{}, "", auth.ErrInvalidCredentials
	}

	token, err := jwtUtils.GenerateToken(user)
	if err != nil {
		loggerVar.Error(auth.ErrGeneratingToken.Error())
		return models.UserResp{}, "", auth.ErrGeneratingToken
	}

	loggerVar.Info("Successful")
	return models.UserResp{Id: user.Id, Login: user.Login}, token, nil
}

func (uc *AuthUsecase) SignUp(ctx context.Context, data models.UserReq) (models.UserResp, string, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	salt := make([]byte, 8)
	rand.Read(salt)
	hashedPassword := validation.HashPassword(salt, data.Password)

	newUser := models.User{
		Id:           uuid.NewV4(),
		Login:        data.Login,
		PasswordHash: hashedPassword,
	}

	if err := uc.repo.InsertUser(ctx, newUser); err != nil {
		switch err {
		case auth.ErrCreatingUser:
			loggerVar.Error(err.Error())
			return models.UserResp{}, "", auth.ErrCreatingUser
		case auth.ErrUserAlreadyExists:
			loggerVar.Error(err.Error())
			return models.UserResp{}, "", auth.ErrUserAlreadyExists
		}
    	return models.UserResp{}, "", err
	}

	token, err := jwtUtils.GenerateToken(newUser)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.UserResp{}, "", auth.ErrGeneratingToken
	}

	loggerVar.Info("Successful")
	return models.UserResp{Id: newUser.Id, Login: newUser.Login}, token, nil
}
