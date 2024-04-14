package repo

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/redis/go-redis/v9"
)

type BlockerRepo struct {
	db  redis.Client
	cfg config.BlockerConfig
}

func CreateBlockerRepo(db redis.Client, cfg config.BlockerConfig) *BlockerRepo {
	return &BlockerRepo{
		db: db,

		cfg: cfg,
	}
}

func (repo *BlockerRepo) GetLoginAttempts(ctx context.Context, ipAddr string) (int, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

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
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	stringCount, err := repo.db.Get(ctx, ipAddr).Result()
	if err != nil {
		logger.Error(err.Error())

		_, err = repo.db.Set(ctx, ipAddr, 1, repo.cfg.RedisExpirationTime).Result()
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		stringCount = "1"
	}

	count, err := strconv.Atoi(stringCount)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = repo.db.Set(ctx, ipAddr, count+1, repo.cfg.RedisExpirationTime).Result()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
