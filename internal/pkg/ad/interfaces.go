package ad

import (
	"context"
	"errors"

	"github.com/K1tten2005/go_vk_intern/internal/models"
)

var (
	ErrCreatingAd = errors.New("ad creation error")
)

type AdUsecase interface {
	CreateAd(ctx context.Context, ad models.Ad) (models.Ad, error)
	GetAds(ctx context.Context, filter models.Filter) ([]models.Ad, error)
}

type AdRepo interface {
	InsertAd(ctx context.Context, ad models.Ad) error
	SelectAds(ctx context.Context, filter models.Filter) ([]models.Ad, error)
}
