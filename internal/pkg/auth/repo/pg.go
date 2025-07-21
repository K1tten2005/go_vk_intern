package repo

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/auth"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/utils/logger"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype/pgxtype"
)

type AuthRepo struct {
	db pgxtype.Querier
}

func CreateAuthRepo(db pgxtype.Querier) *AuthRepo {
	return &AuthRepo{db: db}
}

//go:embed sql/insertUser.sql
var insertUser string

//go:embed sql/selectUserByLogin.sql
var selectUserByLogin string

func (repo *AuthRepo) InsertUser(ctx context.Context, user models.User) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertUser, user.Id, user.Login, user.PasswordHash)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" { 
			return auth.ErrUserAlreadyExists
		}
	}
	if err != nil {
		loggerVar.Error(err.Error())
		return auth.ErrCreatingUser
	}
	loggerVar.Info("Successful")
	return nil
}

func (repo *AuthRepo) SelectUserByLogin(ctx context.Context, login string) (models.User, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	resultUser := models.User{Login: login}
	if err := repo.db.QueryRow(ctx, selectUserByLogin, login).Scan(
		&resultUser.Id,
		&resultUser.PasswordHash,
	); err != nil {
		loggerVar.Error(err.Error())
		return models.User{}, err
	}
	resultUser.Sanitize()

	loggerVar.Info("Successful")
	return resultUser, nil
}
