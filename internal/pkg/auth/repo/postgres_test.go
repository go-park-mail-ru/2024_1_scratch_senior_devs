package repo

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/request"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestAuthRepo_CreateUser(t *testing.T) {
	userId := uuid.NewV4()
	currTime := time.Now().UTC()

	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool)
		err            error
	}{
		{
			name: "CreateUser_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(gomock.Any(), createUser,
					userId,
					"",
					"testuser2",
					request.GetHash("f28fhc2o4m3"),
					currTime,
					"default.jpg",
					sql.NullString{},
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

			repo := CreateAuthRepo(mockPool, testLogger)
			err := repo.CreateUser(context.Background(), models.User{
				Id:           userId,
				Description:  "",
				Username:     "testuser2",
				PasswordHash: request.GetHash("f28fhc2o4m3"),
				CreateTime:   currTime,
				ImagePath:    "default.jpg",
				SecondFactor: "",
			})

			assert.Equal(t, tt.err, err)
		})
	}
}

func TestAuthRepo_GetUserById(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows, uuid.UUID)
		userId         uuid.UUID
		columns        []string
		expectedErr    error
	}{
		{
			name: "GetUserById_Success_1",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserById, id).Return(pgxRows)
				pgxRows.Next()
			},
			userId:      uuid.NewV4(),
			columns:     []string{"description", "username", "password_hash", "create_time", "image_path", "secret"},
			expectedErr: nil,
		},
		{
			name: "GetUserById_Success_2",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserById, id).Return(pgxRows)
				pgxRows.Next()
			},
			userId:      uuid.NewV4(),
			columns:     []string{"description", "username", "password_hash", "create_time", "image_path", "secret"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			description := sql.NullString{}
			secret := sql.NullString{}
			if tt.name == "GetUserById_Success_2" {
				description = sql.NullString{String: "description", Valid: true}
				secret = sql.NullString{String: "fnwovnw", Valid: true}
			}
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(description, "", "", time.Now(), "", secret).ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows, tt.userId)

			repo := CreateAuthRepo(mockPool, testLogger)
			_, err := repo.GetUserById(context.Background(), tt.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAuthRepo_GetUserByUsername(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, pgx.Rows, string)
		username       string
		columns        []string
		expectedErr    error
	}{
		{
			name: "GetUserByUsername_Success_1",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserByUsername, username).Return(pgxRows)
				pgxRows.Next()
			},
			username:    "testuser",
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "secret"},
			expectedErr: nil,
		},
		{
			name: "GetUserByUsername_Success_2",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserByUsername, username).Return(pgxRows)
				pgxRows.Next()
			},
			username:    "testuser",
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "secret"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			defer ctrl.Finish()

			description := sql.NullString{}
			secret := sql.NullString{}
			if tt.name == "GetUserByUsername_Success_2" {
				description = sql.NullString{String: "description", Valid: true}
				secret = sql.NullString{String: "fnwovnw", Valid: true}
			}
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), description, "", time.Now(), "", secret).ToPgxRows()

			tt.mockRepoAction(mockPool, pgxRows, tt.username)

			repo := CreateAuthRepo(mockPool, testLogger)
			_, err := repo.GetUserByUsername(context.Background(), tt.username)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
