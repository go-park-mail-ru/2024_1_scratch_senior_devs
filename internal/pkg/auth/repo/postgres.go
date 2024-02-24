package repo

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

const (
	createUser        = "INSERT INTO users(id, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5);"
	getUserById       = "SELECT (username, password_hash, create_time, image_path) FROM users WHERE id = $1;"
	getUserByUsername = "SELECT (id, password_hash, create_time, image_path) FROM users WHERE username = $1;"
)

type AuthRepo struct {
	db pgxtype.Querier
}

func CreateAuthRepo(db pgxtype.Querier) *AuthRepo {
	return &AuthRepo{db: db}
}

func (repo *AuthRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.Exec(ctx, createUser, user.Id, user.Username, user.PasswordHash, user.CreateTime, user.ImagePath)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (repo *AuthRepo) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	resultUser := &models.User{Id: id}

	err := repo.db.QueryRow(ctx, getUserById, id).Scan(
		&resultUser.Username,
		&resultUser.PasswordHash,
		&resultUser.CreateTime,
		&resultUser.ImagePath,
	)

	if err != nil {
		return &models.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return resultUser, nil
}

func (repo *AuthRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	resultUser := &models.User{Username: username}

	err := repo.db.QueryRow(ctx, getUserByUsername, username).Scan(
		&resultUser.Id,
		&resultUser.PasswordHash,
		&resultUser.CreateTime,
		&resultUser.ImagePath,
	)

	if err != nil {
		return &models.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return resultUser, nil
}
