package repo

import (
	"context"
	_ "embed"
	"log/slog"

	advertisement "github.com/K1tten2005/go_vk_intern/internal/pkg/ad"
	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

type AdRepo struct {
	db pgxtype.Querier
}

func CreateAdRepo(db pgxtype.Querier) *AdRepo {
	return &AdRepo{db: db}
}

//go:embed sql/insertAd.sql
var insertAd string

func (repo *AdRepo) InsertAd(ctx context.Context, ad models.Ad, userId uuid.UUID) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertAd, ad.Id, userId, ad.Title, ad.Description, ad.Price, ad.ImageURL, ad.CreatedAt)
	if err != nil {
		loggerVar.Error(err.Error())
		return advertisement.ErrCreatingAd
	}
	loggerVar.Info("Successful")
	return nil
}
