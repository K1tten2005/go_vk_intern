package repo

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	advertisement "github.com/K1tten2005/go_vk_intern/internal/pkg/ad"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/jackc/pgtype/pgxtype"
)

type AdRepo struct {
	db pgxtype.Querier
}

func CreateAdRepo(db pgxtype.Querier) *AdRepo {
	return &AdRepo{db: db}
}

//go:embed sql/insertAd.sql
var insertAd string

//go:embed sql/selectAds.sql
var selectAds string

func (repo *AdRepo) InsertAd(ctx context.Context, ad models.Ad) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertAd, ad.Id, ad.UserId, ad.Title, ad.Description, ad.Price, ad.ImageURL, ad.CreatedAt)
	if err != nil {
		loggerVar.Error(err.Error())
		return advertisement.ErrCreatingAd
	}
	loggerVar.Info("Successful")
	return nil
}

func (r *AdRepo) SelectAds(ctx context.Context, filter models.Filter) ([]models.Ad, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	query := fmt.Sprintf(selectAds, filter.SortBy, filter.Order)

	offset := (filter.Page - 1) * filter.Limit
	rows, err := r.db.Query(ctx, query, filter.PriceMin, filter.PriceMax, filter.Limit, offset)
	if err != nil {
		loggerVar.Error("query error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		if err := rows.Scan(&ad.Id, &ad.UserId, &ad.Title, &ad.Description, &ad.Price, &ad.ImageURL, &ad.CreatedAt, &ad.AuthorLogin); err != nil {
			loggerVar.Error("scan error: " + err.Error())
			continue
		}
		ads = append(ads, ad)
	}
	return ads, nil
}
