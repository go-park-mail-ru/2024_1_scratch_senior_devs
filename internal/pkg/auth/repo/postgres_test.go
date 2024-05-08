package repo

import (
	"context"
	"database/sql"
	"errors"
	mock_metrics "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics/mocks"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
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
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics)
		err            error
	}{
		{
			name: "CreateUser_Success",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics) {
				mockPool.EXPECT().Exec(gomock.Any(), createUser,
					userId,
					"",
					"testuser2",
					responses.GetHash("f28fhc2o4m3"),
					currTime,
					"default.jpg",
					sql.NullString{},
				).Return(nil, nil)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			tt.mockRepoAction(mockPool, mockMetrics)

			repo := CreateAuthRepo(mockPool, mockMetrics)
			err := repo.CreateUser(context.Background(), models.User{
				Id:           userId,
				Description:  "",
				Username:     "testuser2",
				PasswordHash: responses.GetHash("f28fhc2o4m3"),
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
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, uuid.UUID)
		userId         uuid.UUID
		columns        []string
		expectedErr    error
	}{
		{
			name: "GetUserById_Success_1",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserById, id).Return(pgxRows)
				pgxRows.Next()
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			userId:      uuid.NewV4(),
			columns:     []string{"description", "username", "password_hash", "create_time", "image_path", "data"},
			expectedErr: nil,
		},
		{
			name: "GetUserById_Success_2",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, id uuid.UUID) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserById, id).Return(pgxRows)
				pgxRows.Next()
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			userId:      uuid.NewV4(),
			columns:     []string{"description", "username", "password_hash", "create_time", "image_path", "data"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			description := sql.NullString{}
			data := sql.NullString{}
			if tt.name == "GetUserById_Success_2" {
				description = sql.NullString{String: "description", Valid: true}
				data = sql.NullString{String: "fnwovnw", Valid: true}
			}
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(description, "", "", time.Now(), "", data).ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.userId)

			repo := CreateAuthRepo(mockPool, mockMetrics)
			_, err := repo.GetUserById(context.Background(), tt.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAuthRepo_GetUserByUsername(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoAction func(*pgxpoolmock.MockPgxPool, *mock_metrics.MockDBMetrics, pgx.Rows, string)
		username       string
		columns        []string
		expectedErr    error
	}{
		{
			name: "GetUserByUsername_Success_1",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserByUsername, username).Return(pgxRows)
				pgxRows.Next()
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			username:    "testuser",
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
			expectedErr: nil,
		},
		{
			name: "GetUserByUsername_Success_2",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().QueryRow(gomock.Any(), getUserByUsername, username).Return(pgxRows)
				pgxRows.Next()
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			username:    "testuser",
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()

			description := sql.NullString{}
			data := sql.NullString{}
			if tt.name == "GetUserByUsername_Success_2" {
				description = sql.NullString{String: "description", Valid: true}
				data = sql.NullString{String: "fnwovnw", Valid: true}
			}
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), description, "", time.Now(), "", data).ToPgxRows()

			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.username)

			repo := CreateAuthRepo(mockPool, mockMetrics)
			_, err := repo.GetUserByUsername(context.Background(), tt.username)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAuthRepo_Deletedata(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name           string
		mockRepoAction func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string)
		args           args
		expectedErr    error
		columns        []string
	}{
		{
			name: "TestSuccess",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().Exec(context.Background(), deleteSecondFactor, "user").Return(nil, nil).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			args: args{
				username: "user",
			},
			expectedErr: nil,
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
		{
			name: "TestFail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().Exec(context.Background(), deleteSecondFactor, "user").Return(nil, errors.New("error")).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			args: args{
				username: "user",
			},
			expectedErr: errors.New("error"),
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), "description", "", time.Now(), "", "data").ToPgxRows()

			repo := CreateAuthRepo(mockPool, mockMetrics)
			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args.username)
			err := repo.DeleteSecret(context.Background(), tt.args.username)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAuthRepo_Updatedata(t *testing.T) {
	type args struct {
		username string
		data     string
	}
	tests := []struct {
		name           string
		mockRepoAction func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string)
		args           args
		expectedErr    error
		columns        []string
	}{
		{
			name: "TestSuccess",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().Exec(context.Background(), updateSecondFactor, "data", "user").Return(nil, nil).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			args: args{
				username: "user",
				data:     "data",
			},
			expectedErr: nil,
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
		{
			name: "TestFail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().Exec(context.Background(), updateSecondFactor, "data", "user").Return(nil, errors.New("error")).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			args: args{
				username: "user",
				data:     "data",
			},
			expectedErr: errors.New("error"),
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(uuid.NewV4(), "description", "", time.Now(), "", tt.args.data).ToPgxRows()

			repo := CreateAuthRepo(mockPool, mockMetrics)
			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args.username)
			err := repo.UpdateSecret(context.Background(), tt.args.username, tt.args.data)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAuthRepo_UpdateProfileAvatar(t *testing.T) {
	userId := uuid.NewV4()
	type args struct {
		Id   uuid.UUID
		path string
	}
	tests := []struct {
		name           string
		mockRepoAction func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string)
		args           args
		expectedErr    error
		columns        []string
	}{
		{
			name: "TestSuccess",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().Exec(context.Background(), updateProfileAvatar, "path", userId).Return(nil, nil).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			args: args{
				Id:   userId,
				path: "path",
			},
			expectedErr: nil,
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
		{
			name: "TestFail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, username string) {
				mockPool.EXPECT().Exec(context.Background(), updateProfileAvatar, "path", userId).Return(nil, errors.New("error")).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			args: args{
				Id:   userId,
				path: "path",
			},
			expectedErr: errors.New("error"),
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(userId, "description", "", time.Now(), "", "").ToPgxRows()

			repo := CreateAuthRepo(mockPool, mockMetrics)
			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args.Id.String())
			err := repo.UpdateProfileAvatar(context.Background(), tt.args.Id, tt.args.path)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAuthRepo_UpdateProfile(t *testing.T) {
	type args struct {
		user models.User
	}
	tests := []struct {
		name           string
		mockRepoAction func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, user models.User)
		args           args
		columns        []string
		expectedErr    error
	}{
		{
			name: "TestSuccess",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, user models.User) {
				mockPool.EXPECT().Exec(context.Background(), updateProfile, user.Description, user.PasswordHash, user.Id).Return(nil, nil).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
			},
			args:        args{},
			expectedErr: nil,
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
		{
			name: "TestFail",
			mockRepoAction: func(mockPool *pgxpoolmock.MockPgxPool, metr *mock_metrics.MockDBMetrics, pgxRows pgx.Rows, user models.User) {
				mockPool.EXPECT().Exec(context.Background(), updateProfile, user.Description, user.PasswordHash, user.Id).Return(nil, errors.New("error")).Times(1)
				metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any()).Return()
				metr.EXPECT().IncreaseErrors(gomock.Any()).Return()
			},
			args:        args{},
			expectedErr: errors.New("error"),
			columns:     []string{"id", "description", "password_hash", "create_time", "image_path", "data"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			mockMetrics := mock_metrics.NewMockDBMetrics(ctrl)
			defer ctrl.Finish()
			pgxRows := pgxpoolmock.NewRows(tt.columns).AddRow(tt.args.user.Id, "description", "", time.Now(), "", "").ToPgxRows()

			repo := CreateAuthRepo(mockPool, mockMetrics)
			tt.mockRepoAction(mockPool, mockMetrics, pgxRows, tt.args.user)
			err := repo.UpdateProfile(context.Background(), tt.args.user)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
