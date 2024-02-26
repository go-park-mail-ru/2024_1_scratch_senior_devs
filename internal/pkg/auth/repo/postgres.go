package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

const (
	createUser        = "INSERT INTO users (id, description, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5, $6);"
	getUserById       = "SELECT description, username, password_hash, create_time, image_path FROM users WHERE id = $1;"
	getUserByUsername = "SELECT id, description, password_hash, create_time, image_path FROM users WHERE username = $1;"
	//getPasswordByUsername = "SELECT (password_hash) FROM users WHERE username = $1;"
)

type AuthRepo struct {
	db pgxtype.Querier
}

func CreateAuthRepo(db pgxtype.Querier) *AuthRepo {
	return &AuthRepo{db: db}
}

func (repo *AuthRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.Exec(ctx, createUser, user.Id, user.Description, user.Username, user.PasswordHash, user.CreateTime, user.ImagePath)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (repo *AuthRepo) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	resultUser := &models.User{Id: id}
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
	} else if err != nil {
		return &models.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return resultUser, nil
}

func (repo *AuthRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	resultUser := &models.User{Username: username}
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
	} else if err != nil {
		return &models.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return resultUser, nil
}

// returns user model if credentials are fine and returns error if not
func (repo *AuthRepo) CheckUserCredentials(ctx context.Context, username string, password string) (*models.User, error) {
	user, err := repo.GetUserByUsername(ctx, username)
	if err != nil {
		return &models.User{}, err
	}
	if user.PasswordHash != password {
		return &models.User{}, fmt.Errorf("wrong credentials")
	}
	return user, nil

}
