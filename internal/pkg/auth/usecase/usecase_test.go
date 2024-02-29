package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/satori/uuid"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	mockAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"

	"github.com/golang/mock/gomock"
)

func TestAuthUsecase_SignUp(t *testing.T) { //тут можем чекнуть по сути только наличие или отсутствие ошибки. плохой тест какой-то

	type args struct {
		data *models.UserFormData
	}
	tests := []struct {
		name       string
		repoMocker func(context.Context, *mockAuth.MockAuthRepo)
		args       args
		//want       *models.User
		///want1      string
		///want2      time.Time
		wantErr bool
	}{
		{
			name: "TestSuccess",
			repoMocker: func(ctx context.Context, repo *mockAuth.MockAuthRepo) {
				repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil).Times(1)
			},
			args: args{
				data: &models.UserFormData{
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
				data: &models.UserFormData{
					Username: "hello",
					Password: "qwerty111",
				},
			},

			wantErr: true,
		},
	}
	/*userData := models.User{
		Id:           uuid.FromStringOrNil("12f6a194-1c9e-4726-b295-53cb0d0bd457"),
		Description:  nil,
		Username:     "hello",
		PasswordHash: getHash("qwerty111"),
		CreateTime:   time.Now().UTC(),
		ImagePath:    "default.jpg",
	}*/

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()
			repo := mockAuth.NewMockAuthRepo(ctl)
			uc := CreateAuthUsecase(repo)

			tt.repoMocker(context.Background(), repo)
			_, _, _, err := uc.SignUp(context.Background(), tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUsecase.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/*if !reflect.DeepEqual(userModel, tt.want) {
				t.Errorf("AuthUsecase.SignUp() got = %v, want %v", userModel, tt.want)
			}
			if token != tt.want1 {
				t.Errorf("AuthUsecase.SignUp() got1 = %v, want %v", token, tt.want1)
			}
			if !reflect.DeepEqual(creationTime, tt.want2) {
				t.Errorf("AuthUsecase.SignUp() got2 = %v, want %v", creationTime, tt.want2)
			}*/
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
		data *models.UserFormData
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
				repo.EXPECT().CheckUserCredentials(gomock.Any(), username, passwordHash).Return(&models.User{
					Id:           uuid.NewV4(),
					Description:  nil,
					Username:     username,
					PasswordHash: passwordHash,
				}, getErr(wantErr)).Times(1)
			},
			args: args{
				data: &models.UserFormData{
					Username: "hello",
					Password: "qwerty111",
				},
			},
			wantErr: false,
		},
		{
			name: "AuthUsecase_SignIn_Fail",
			repoMocker: func(repo *mockAuth.MockAuthRepo, username string, passwordHash string, wantErr bool) {
				repo.EXPECT().CheckUserCredentials(gomock.Any(), username, passwordHash).Return(&models.User{
					Id:           uuid.NewV4(),
					Description:  nil,
					Username:     username,
					PasswordHash: passwordHash,
				}, getErr(wantErr)).Times(1)
			},
			args: args{
				data: &models.UserFormData{
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
			uc := CreateAuthUsecase(repo)
			defer ctrl.Finish()

			tt.repoMocker(repo, tt.args.data.Username, utils.GetHash(tt.args.data.Password), tt.wantErr)
			_, _, _, err := uc.SignIn(context.Background(), tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUsecase.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
