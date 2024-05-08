package repo

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/redis/go-redis/v9"
)

type BlockerRepo struct {
	db   redis.Client
	cfg  config.BlockerConfig
	metr metrics.DBMetrics
}

func CreateBlockerRepo(db redis.Client, cfg config.BlockerConfig, metr metrics.DBMetrics) *BlockerRepo {
	return &BlockerRepo{
		db:   db,
		cfg:  cfg,
		metr: metr,
	}
}

func (repo *BlockerRepo) GetLoginAttempts(ctx context.Context, ipAddr string) (int, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	stringCount, err := repo.db.Get(ctx, ipAddr).Result()
	repo.metr.ObserveResponseTime("getByKey", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("getByKey")
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

	start := time.Now()
	stringCount, err := repo.db.Get(ctx, ipAddr).Result()
	repo.metr.ObserveResponseTime("getByKey", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())

		start = time.Now()
		_, err = repo.db.Set(ctx, ipAddr, 1, repo.cfg.RedisExpirationTime).Result()
		repo.metr.ObserveResponseTime("setByKey", time.Since(start).Seconds())
		if err != nil {
			logger.Error(err.Error())
			repo.metr.IncreaseErrors("setByKey")
			return err
		}
		stringCount = "1"
	}

	count, err := strconv.Atoi(stringCount)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	start = time.Now()
	_, err = repo.db.Set(ctx, ipAddr, count+1, repo.cfg.RedisExpirationTime).Result()
	repo.metr.ObserveResponseTime("setByKey", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("setByKey")
		return err
	}

	logger.Info("success")
	return nil
}
