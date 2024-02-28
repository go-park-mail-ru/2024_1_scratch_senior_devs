package repo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

func TestAuthRepo_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "CreateUser_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), createUser,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, nil)
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool)

			repo := CreateAuthRepo(mockPool)
			err := repo.CreateUser(context.Background(), &models.User{})

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAuthRepo_GetUserById(t *testing.T) {
	userId := uuid.NewV4()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows)
		columns        []string
		err            error
	}{
		{
			name: "GetUserById_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserById, userId).Return(pgxRows)
				pgxRows.Next()
			},
			columns: []string{"description", "username", "password_hash", "create_time", "image_path"},
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(sql.NullString{}, "", "", time.Now(), "").ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows)

			repo := CreateAuthRepo(mockPool)
			_, err := repo.GetUserById(context.Background(), userId)

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAuthRepo_GetUserByUsername(t *testing.T) {
	username := "test_user"

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows)
		columns        []string
		err            error
	}{
		{
			name: "GetUserByUsername_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserByUsername, username).Return(pgxRows)
				pgxRows.Next()
			},
			columns: []string{"id", "description", "password_hash", "create_time", "image_path"},
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), sql.NullString{}, "", time.Now(), "").ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows)

			repo := CreateAuthRepo(mockPool)
			_, err := repo.GetUserByUsername(context.Background(), username)

			assert.Equal(t, tt.err, err)
		})
	}
}
