package ad

import (
	"context"
	"errors"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/satori/uuid"
)

var (
	ErrCreatingAd = errors.New("ad creation error")
)

type AdUsecase interface {
	CreateAd(ctx context.Context, ad models.Ad, userId uuid.UUID) (models.Ad, error)
}

type AdRepo interface {
	InsertAd(ctx context.Context, ad models.Ad, userId uuid.UUID) error
}
