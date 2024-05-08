package repo

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"log/slog"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	createUser          = "INSERT INTO users(id, description, username, password_hash, create_time, image_path, secret) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	getUserById         = "SELECT description, username, password_hash, create_time, image_path, secret FROM users WHERE id = $1;"
	getUserByUsername   = "SELECT id, description, password_hash, create_time, image_path, secret FROM users WHERE username = $1;"
	updateProfile       = "UPDATE users SET description = $1, password_hash = $2 WHERE id = $3;"
	updateProfileAvatar = "UPDATE users SET image_path = $1 WHERE id = $2;"
	updateSecondFactor  = "UPDATE users SET secret = $1 WHERE username = $2;"
	deleteSecondFactor  = "UPDATE users SET secret = NULL WHERE username = $1;"
)

type AuthRepo struct {
	db   pgxtype.Querier
	metr metrics.DBMetrics
}

func CreateAuthRepo(db pgxtype.Querier, metr metrics.DBMetrics) *AuthRepo {
	return &AuthRepo{
		db:   db,
		metr: metr,
	}
}

func (repo *AuthRepo) CreateUser(ctx context.Context, user models.User) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	userSecret := sql.NullString{
		String: string(user.SecondFactor),
		Valid:  user.SecondFactor != "",
	}

	start := time.Now()
	_, err := repo.db.Exec(ctx, createUser, user.Id, user.Description, user.Username, user.PasswordHash, user.CreateTime, user.ImagePath, userSecret)
	repo.metr.ObserveResponseTime("createUser", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("createUser")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *AuthRepo) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultUser := models.User{Id: id}
	description := sql.NullString{}
	secret := sql.NullString{}

	start := time.Now()
	err := repo.db.QueryRow(ctx, getUserById, id).Scan(
		&description,
		&resultUser.Username,
		&resultUser.PasswordHash,
		&resultUser.CreateTime,
		&resultUser.ImagePath,
		&secret,
	)
	repo.metr.ObserveResponseTime("getUserById", time.Since(start).Seconds())

	if description.Valid {
		resultUser.Description = description.String
	}

	if secret.Valid {
		resultUser.SecondFactor = models.Secret(secret.String)
	}

	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("getUserById")
		return models.User{}, err
	}

	logger.Info("success")
	return resultUser, nil
}

func (repo *AuthRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultUser := models.User{Username: username}
	description := sql.NullString{}
	secret := sql.NullString{}

	start := time.Now()
	err := repo.db.QueryRow(ctx, getUserByUsername, username).Scan(
		&resultUser.Id,
		&description,
		&resultUser.PasswordHash,
		&resultUser.CreateTime,
		&resultUser.ImagePath,
		&secret,
	)
	repo.metr.ObserveResponseTime("getUserByUsername", time.Since(start).Seconds())

	if description.Valid {
		resultUser.Description = description.String
	}

	if secret.Valid {
		resultUser.SecondFactor = models.Secret(secret.String)
	}

	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("getUserByUsername")
		return models.User{}, err
	}

	logger.Info("success")
	return resultUser, nil
}

func (repo *AuthRepo) UpdateProfile(ctx context.Context, user models.User) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, updateProfile, user.Description, user.PasswordHash, user.Id)
	repo.metr.ObserveResponseTime("updateProfile", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("updateProfile")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *AuthRepo) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, imagePath string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, updateProfileAvatar, imagePath, userID)
	repo.metr.ObserveResponseTime("updateProfileAvatar", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("updateProfileAvatar")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *AuthRepo) UpdateSecret(ctx context.Context, username string, secret string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, updateSecondFactor, secret, username)
	repo.metr.ObserveResponseTime("updateSecondFactor", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("updateSecondFactor")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *AuthRepo) DeleteSecret(ctx context.Context, username string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, deleteSecondFactor, username)
	repo.metr.ObserveResponseTime("deleteSecondFactor", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("deleteSecondFactor")
		return err
	}

	logger.Info("success")
	return nil
}
