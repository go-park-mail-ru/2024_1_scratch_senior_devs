package repo

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
)

const (
	createUser          = "INSERT INTO users(id, description, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5, $6);"
	getUserById         = "SELECT description, username, password_hash, create_time, image_path FROM users WHERE id = $1;"
	getUserByUsername   = "SELECT id, description, password_hash, create_time, image_path FROM users WHERE username = $1;"
	updateProfile       = "UPDATE users SET description = $1, password_hash = $2 WHERE id = $3;"
	updateProfileAvatar = "UPDATE users SET image_path = $1 WHERE id = $2;"
)

type AuthRepo struct {
	db     pgxtype.Querier
	logger *slog.Logger
}

func CreateAuthRepo(db pgxtype.Querier, logger *slog.Logger) *AuthRepo {
	return &AuthRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *AuthRepo) CreateUser(ctx context.Context, user models.User) error {
	logger := repo.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	_, err := repo.db.Exec(ctx, createUser, user.Id, user.Description, user.Username, user.PasswordHash, user.CreateTime, user.ImagePath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *AuthRepo) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	logger := repo.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	resultUser := models.User{Id: id}
	description := sql.NullString{}

	err := repo.db.QueryRow(ctx, getUserById, id).Scan(
		&description,
		&resultUser.Username,
		&resultUser.PasswordHash,
		&resultUser.CreateTime,
		&resultUser.ImagePath,
	)

	if description.Valid {
		resultUser.Description = description.String
	}

	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	logger.Info("success")
	return resultUser, nil
}

func (repo *AuthRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	logger := repo.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	resultUser := models.User{Username: username}
	description := sql.NullString{}

	err := repo.db.QueryRow(ctx, getUserByUsername, username).Scan(
		&resultUser.Id,
		&description,
		&resultUser.PasswordHash,
		&resultUser.CreateTime,
		&resultUser.ImagePath,
	)

	if description.Valid {
		resultUser.Description = description.String
	}

	if err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	logger.Info("success")
	return resultUser, nil
}

func (repo *AuthRepo) UpdateProfile(ctx context.Context, user models.User) error {
	logger := repo.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	_, err := repo.db.Exec(ctx, updateProfile, user.Description, user.PasswordHash, user.Id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *AuthRepo) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, imagePath string) error {
	logger := repo.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	_, err := repo.db.Exec(ctx, updateProfileAvatar, imagePath, userID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
