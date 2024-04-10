package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	mock_note "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/mocks"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mockAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"

	"github.com/golang/mock/gomock"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestAuthUsecase_SignUp(t *testing.T) {
	testConfig := config.Config{
		AuthUsecase: config.AuthUsecaseConfig{
			DefaultImagePath: "default.jpg",
			JWTLifeTime:      time.Duration(24 * time.Hour),
		},
		UserValidation: config.UserValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
	type args struct {
		data models.UserFormData
	}
	tests := []struct {
		name       string
		repoMocker func(context.Context, *mockAuth.MockAuthRepo, *mock_note.MockNoteRepo)
		args       args
		wantErr    bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mockAuth.MockAuthRepo, noteRepo *mock_note.MockNoteRepo) {
				repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil).Times(1)
				noteRepo.EXPECT().CreateNote(ctx, gomock.Any()).Return(nil).Times(1)
				noteRepo.EXPECT().MakeHelloNoteData(gomock.Any()).Return([]byte{}).Times(1)
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
			name: "TestFail",
			repoMocker: func(ctx context.Context, repo *mockAuth.MockAuthRepo, noteRepo *mock_note.MockNoteRepo) {
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
			noteRepo := mock_note.NewMockNoteRepo(ctl)
			uc := CreateAuthUsecase(repo, noteRepo, testLogger, testConfig.AuthUsecase, testConfig.UserValidation)

			tt.repoMocker(context.Background(), repo, noteRepo)
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
	testConfig := config.Config{
		AuthUsecase: config.AuthUsecaseConfig{
			DefaultImagePath: "default.jpg",
			JWTLifeTime:      time.Duration(24 * time.Hour),
		},
		UserValidation: config.UserValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
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
			noteRepo := mock_note.NewMockNoteRepo(ctrl)
			uc := CreateAuthUsecase(repo, noteRepo, testLogger, testConfig.AuthUsecase, testConfig.UserValidation)
			defer ctrl.Finish()

			tt.repoMocker(repo, tt.args.data.Username, responses.GetHash(tt.args.data.Password), tt.wantErr)
			_, _, _, err := uc.SignIn(context.Background(), tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUsecase.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCheckUser(t *testing.T) {
	testConfig := config.Config{
		AuthUsecase: config.AuthUsecaseConfig{
			DefaultImagePath: "default.jpg",
			JWTLifeTime:      time.Duration(24 * time.Hour),
		},
		UserValidation: config.UserValidationConfig{
			MinUsernameLength:    4,
			MaxUsernameLength:    12,
			MinPasswordLength:    8,
			MaxPasswordLength:    20,
			PasswordAllowedExtra: "$%&#",
			SecretLength:         6,
		},
	}
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
					PasswordHash: responses.GetHash("f34ovin332"),
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
					PasswordHash: responses.GetHash("f34ovin332"),
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
			noteRepo := mock_note.NewMockNoteRepo(ctrl)
			uc := CreateAuthUsecase(repo, noteRepo, testLogger, testConfig.AuthUsecase, testConfig.UserValidation)
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
