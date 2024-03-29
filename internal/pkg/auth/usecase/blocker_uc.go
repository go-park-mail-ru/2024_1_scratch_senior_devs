package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
)

type BlockerUsecase struct {
	repo   auth.BlockerRepo
	logger *slog.Logger
}

func CreateBlockerUsecase(repo auth.BlockerRepo, logger *slog.Logger) *BlockerUsecase {
	return &BlockerUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (bu *BlockerUsecase) CheckLoginAttempts(ctx context.Context, ipAddress string) error {
	err := bu.repo.IncreaseLoginAttempts(ctx, ipAddress)
	if err != nil {
		return err
	}

	requestsMade, err := bu.repo.GetLoginAttempts(ctx, ipAddress)
	if err != nil {
		return err
	}

	if requestsMade > config.MaxWrongRequests {
		return errors.New("too many attempts")
	}
	return nil
}