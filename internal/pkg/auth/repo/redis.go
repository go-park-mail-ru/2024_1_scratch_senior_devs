package repo

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const RedisExpirationTime = time.Minute

type BlockerRepo struct {
	db     redis.Client
	logger *slog.Logger
}

func CreateBlockerRepo(db redis.Client, logger *slog.Logger) *BlockerRepo {
	return &BlockerRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *BlockerRepo) GetLoginAttempts(ctx context.Context, ipAddr string) (int, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	stringCount, err := repo.db.Get(ctx, ipAddr).Result()
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	count, err := strconv.Atoi(stringCount)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	logger.Info("success")
	return count, nil
}

func (repo *BlockerRepo) IncreaseLoginAttempts(ctx context.Context, ipAddr string) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	stringCount, err := repo.db.Get(ctx, ipAddr).Result()
	if err != nil {
		logger.Error(err.Error())

		_, err = repo.db.Set(ctx, ipAddr, 1, RedisExpirationTime).Result()
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	count, err := strconv.Atoi(stringCount)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = repo.db.Set(ctx, ipAddr, count+1, RedisExpirationTime).Result()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
