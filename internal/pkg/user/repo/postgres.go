package repo

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	createUser        = "INSERT INTO users(id, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5);"
	getUserById       = "SELECT (username, password_hash, create_time, image_path) FROM users WHERE id = $1"
	getUserByUsername = "SELECT (id, password_hash, create_time, image_path) FROM users WHERE username = $1"
)

type UsersRepo struct {
	db pgxtype.Querier
}

func CreateProfileRepo(db pgxtype.Querier) *UsersRepo {
	return &UsersRepo{db: db}
}

func (repo *UsersRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.Exec(ctx, createUser, user.Id, user.Username, user.PasswordHash, user.CreateTime, user.ImagePath)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (repo *UsersRepo) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
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

func (repo *UsersRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
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
