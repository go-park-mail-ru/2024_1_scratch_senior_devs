package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/delivery"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mockAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"

	"github.com/golang/mock/gomock"
)

var testLogger *slog.Logger
var testConfig *config.Config

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	testConfig = config.LoadConfig("../../config/config.yaml", testLogger)
}

func TestAuthUsecase_SignUp(t *testing.T) { //тут можем чекнуть по сути только наличие или отсутствие ошибки. плохой тест какой-то

	type args struct {
		data models.UserFormData
	}
	tests := []struct {
		name       string
		repoMocker func(context.Context, *mockAuth.MockAuthRepo)
		args       args
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mockAuth.MockAuthRepo) {
				repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				data: models.UserFormData{
					Username: "hello",
					Password: "qwerty111",
				},
			},
			wantErr: false,
		},
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mockAuth.MockAuthRepo) {
				repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(errors.New("error creating user")).Times(1)
			},
			args: args{
				data: models.UserFormData{
					Username: "hello",
					Password: "qwerty111",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mockAuth.NewMockAuthRepo(ctl)
			uc := CreateAuthUsecase(repo, testLogger, testConfig.AuthUsecase, testConfig.UserValidation)

			tt.repoMocker(context.Background(), repo)
			_, _, _, err := uc.SignUp(context.Background(), tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUsecase.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func getErr(wantErr bool) error {
	if wantErr {
		return errors.New("error")
	}
	return nil
}

func TestAuthUsecase_SignIn(t *testing.T) {
	type args struct {
		data models.UserFormData
	}
	tests := []struct {
		name       string
		repoMocker func(*mockAuth.MockAuthRepo, string, string, bool)
		args       args
		wantErr    bool
	}{
		{
			name: "AuthUsecase_SignIn_Success",
			repoMocker: func(repo *mockAuth.MockAuthRepo, username string, passwordHash string, wantErr bool) {
				repo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(models.User{
					Id:           uuid.NewV4(),
					Description:  "",
					Username:     username,
					PasswordHash: passwordHash,
				}, getErr(wantErr)).Times(1)
			},
			args: args{
				data: models.UserFormData{
					Username: "hello",
					Password: "qwerty111",
				},
			},
			wantErr: false,
		},
		{
			name: "AuthUsecase_SignIn_Fail",
			repoMocker: func(repo *mockAuth.MockAuthRepo, username string, passwordHash string, wantErr bool) {
				repo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(models.User{
					Id:           uuid.NewV4(),
					Description:  "",
					Username:     username,
					PasswordHash: passwordHash,
				}, getErr(wantErr)).Times(1)
			},
			args: args{
				data: models.UserFormData{
					Username: "hello",
					Password: "qwerty111",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mockAuth.NewMockAuthRepo(ctrl)
			uc := CreateAuthUsecase(repo, testLogger, testConfig.AuthUsecase, testConfig.UserValidation)
			defer ctrl.Finish()

			tt.repoMocker(repo, tt.args.data.Username, delivery.GetHash(tt.args.data.Password), tt.wantErr)
			_, _, _, err := uc.SignIn(context.Background(), tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUsecase.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCheckUser(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name       string
		repoMocker func(*mockAuth.MockAuthRepo, uuid.UUID, bool)
		args       args
		wantErr    bool
	}{
		{
			name: "AuthUsecase_CheckUser_Success",
			repoMocker: func(repo *mockAuth.MockAuthRepo, id uuid.UUID, wantErr bool) {
				repo.EXPECT().GetUserById(gomock.Any(), id).Return(models.User{
					Id:           uuid.NewV4(),
					Description:  "",
					Username:     "testuser1",
					PasswordHash: delivery.GetHash("f34ovin332"),
				}, getErr(wantErr)).Times(1)
			},
			args: args{
				id: uuid.NewV4(),
			},
			wantErr: false,
		},
		{
			name: "AuthUsecase_CheckUser_Fail",
			repoMocker: func(repo *mockAuth.MockAuthRepo, id uuid.UUID, wantErr bool) {
				repo.EXPECT().GetUserById(gomock.Any(), id).Return(models.User{
					Id:           uuid.NewV4(),
					Description:  "",
					Username:     "testuser1",
					PasswordHash: delivery.GetHash("f34ovin332"),
				}, getErr(wantErr)).Times(1)
			},
			args: args{
				id: uuid.NewV4(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mockAuth.NewMockAuthRepo(ctrl)
			uc := CreateAuthUsecase(repo, testLogger, testConfig.AuthUsecase, testConfig.UserValidation)
			defer ctrl.Finish()

			tt.repoMocker(repo, tt.args.id, tt.wantErr)
			_, err := uc.CheckUser(context.Background(), tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUsecase.CheckUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
