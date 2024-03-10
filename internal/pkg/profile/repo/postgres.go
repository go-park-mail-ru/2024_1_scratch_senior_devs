package repo

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

const (
	getUserById   = "SELECT description, username, password_hash, create_time, image_path FROM users WHERE id = $1;"
	updateProfile = "UPDATE users SET description = $1, passwordHash = $2 WHERE id = $3;"
)

type ProfileRepo struct {
	db pgxtype.Querier
}

func CreateProfileRepo(db pgxtype.Querier) *ProfileRepo {
	return &ProfileRepo{
		db: db,
	}
}

func (repo *ProfileRepo) UpdateProfile(ctx context.Context, user models.User) error {
	_, err := repo.db.Exec(ctx, updateProfile, user.Description, user.PasswordHash, user.Id)
	return err
}

func (repo *ProfileRepo) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
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
		return models.User{}, err
	}

	return resultUser, nil
}
