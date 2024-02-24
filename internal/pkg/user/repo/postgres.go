package repo

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	createUser      = "INSERT INTO users(id, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5);"
	getUserPassword = "SELECT password_hash FROM users WHERE username = $1"
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

func (repo *UsersRepo) GetUser(ctx context.Context, credentials models.JWTPayload) (string, error) {
	var password string
	err := repo.db.QueryRow(ctx, getUserPassword, credentials.Username).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("error getting user's password: %w", err)
	}

	return password, nil
} //ff
