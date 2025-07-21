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

func (uc *AdUsecase) CreateAd(ctx context.Context, data models.Ad) (models.Ad, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	data.CreatedAt = time.Now()
	data.Id = uuid.NewV4()

	if err := uc.repo.InsertAd(ctx, data); err != nil {
		loggerVar.Error(err.Error())
		return models.Ad{}, err
	}

	loggerVar.Info("Successful")
	return data, nil
}

func (uc *AdUsecase) GetAds(ctx context.Context, filter models.Filter) ([]models.Ad, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	ads, err := uc.repo.SelectAds(ctx, filter)
	if err != nil {
		loggerVar.Error("error fetching ads: " + err.Error())
		return ads, err
	}

	return ads, nil
}
