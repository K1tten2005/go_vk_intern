package repo

import (
	"context"
	"testing"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/auth"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	user := models.User{
		Id:           uuid.NewV4(),
		Login:        "test2025",
		PasswordHash: []byte("hashed_password"),
	}

	tests := []struct {
		name        string
		repoMocker  func(*pgxpoolmock.MockPgxPool)
		expectedErr error
	}{
		{
			name: "Success",
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(
					gomock.Any(),
					gomock.Any(),
					user.Id,
					user.Login,
					user.PasswordHash,
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Duplicate login",
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockErr := &pgconn.PgError{Code: "23505"}
				mockPool.EXPECT().Exec(
					gomock.Any(),
					gomock.Any(),
					user.Id,
					user.Login,
					user.PasswordHash,
				).Return(nil, mockErr)
			},
			expectedErr: auth.ErrUserAlreadyExists,
		},
		{
			name: "Other error",
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(
					gomock.Any(),
					gomock.Any(),
					user.Id,
					user.Login,
					user.PasswordHash,
				).Return(nil, assert.AnError)
			},
			expectedErr: auth.ErrCreatingUser,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			test.repoMocker(mockPool)

			authRepo := CreateAuthRepo(mockPool)
			err := authRepo.InsertUser(context.Background(), user)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

