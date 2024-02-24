package repo

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	createUser = "INSERT INTO users(id, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5);"
)

type UsersRepo struct {
	db *sql.DB
}

func CreateProfileRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (repo *UsersRepo) CreateUser(user *models.User) error {
	_, err := repo.db.Exec(createUser, user.Id, user.Username, user.PasswordHash, user.CreateTime, user.ImagePath)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}
