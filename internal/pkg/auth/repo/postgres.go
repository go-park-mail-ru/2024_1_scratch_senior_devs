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
	createUser                = "INSERT INTO users(id, description, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5, $6);"
	getUserById               = "SELECT description, username, password_hash, create_time, image_path FROM users WHERE id = $1;"
	getUserByUsername         = "SELECT id, description, password_hash, create_time, image_path FROM users WHERE username = $1;"
	getPasswordHashByUsername = "SELECT password_hash FROM users WHERE username = $1"
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
	repo.logger.Info(utils.GFN())
	_, err := repo.db.Exec(ctx, createUser, user.Id, user.Description, user.Username, user.PasswordHash, user.CreateTime, user.ImagePath)
	if err != nil {
		repo.logger.Error(err.Error())
	}
	return err
}

func (repo *AuthRepo) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	repo.logger.Info(utils.GFN())
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
		repo.logger.Error(err.Error())
		return models.User{}, err
	}

	return resultUser, nil
}

func (repo *AuthRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
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
		repo.logger.Error(err.Error())
		return models.User{}, err
	}

	return resultUser, nil
}
