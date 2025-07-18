package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/ad"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/satori/uuid"
)

type AdUsecase struct {
	repo ad.AdRepo
}

func CreateAdUsecase(repo ad.AdRepo) *AdUsecase {
	return &AdUsecase{repo: repo}
}

func (uc *AdUsecase) CreateAd(ctx context.Context, data models.Ad, userId uuid.UUID) (models.Ad, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	data.CreatedAt = time.Now()
	data.Id = uuid.NewV4()
	data.Price *= 100

	err := uc.repo.InsertAd(ctx, data, userId)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Ad{}, err
	}

	loggerVar.Info("Successful")
	return data, nil
}
