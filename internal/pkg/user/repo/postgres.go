package repo

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	createUser = "INSERT INTO users(id, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5);"
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
