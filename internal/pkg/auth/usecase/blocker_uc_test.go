package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_auth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBlockerUsecase_CheckLoginAttempts(t *testing.T) {
	testConfig := config.Config{
		Blocker: config.BlockerConfig{
			RedisExpirationTime: time.Duration(time.Minute),
			MaxWrongRequests:    5,
		},
	}
	type args struct {
		ctx       context.Context
		ipAddress string
	}
	tests := []struct {
		name        string
		repoMocker  func(ctx context.Context, repo *mock_auth.MockBlockerRepo)
		args        args
		expectedErr error
	}{
		{
			name: "TestBlockerUC_Success",
			repoMocker: func(ctx context.Context, repo *mock_auth.MockBlockerRepo) {
				repo.EXPECT().IncreaseLoginAttempts(ctx, "1").Return(nil)
				repo.EXPECT().GetLoginAttempts(ctx, "1").Return(1, nil)
			},
			args: args{
				ctx:       context.Background(),
				ipAddress: "1",
			},
			expectedErr: nil,
		},
		{
			name: "TestBlockerUC_Fail_Too_Many_Requests",
			repoMocker: func(ctx context.Context, repo *mock_auth.MockBlockerRepo) {
				repo.EXPECT().IncreaseLoginAttempts(ctx, "1").Return(nil)
				repo.EXPECT().GetLoginAttempts(ctx, "1").Return(10, nil)
			},
			args: args{
				ctx:       context.Background(),
				ipAddress: "1",
			},
			expectedErr: errors.New("too many attempts"),
		},
		{
			name: "TestBlockerUC_Fail_To_Increase",
			repoMocker: func(ctx context.Context, repo *mock_auth.MockBlockerRepo) {
				repo.EXPECT().IncreaseLoginAttempts(ctx, "1").Return(errors.New("can`t increase")).Times(1)
			},
			args: args{
				ctx:       context.Background(),
				ipAddress: "1",
			},
			expectedErr: errors.New("can`t increase"),
		},
		{
			name: "TestBlockerUC_Fail_To_Get_Attempts",
			repoMocker: func(ctx context.Context, repo *mock_auth.MockBlockerRepo) {
				repo.EXPECT().IncreaseLoginAttempts(ctx, "1").Return(nil).Times(1)
				repo.EXPECT().GetLoginAttempts(ctx, "1").Return(1, errors.New("failed to get")).Times(1)
			},
			args: args{
				ctx:       context.Background(),
				ipAddress: "1",
			},
			expectedErr: errors.New("failed to get"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_auth.NewMockBlockerRepo(ctl)
			uc := CreateBlockerUsecase(repo, testConfig.Blocker)

			tt.repoMocker(context.Background(), repo)

			err := uc.CheckLoginAttempts(tt.args.ctx, tt.args.ipAddress)
			assert.Equal(t, tt.expectedErr, err)

		})
	}
}
