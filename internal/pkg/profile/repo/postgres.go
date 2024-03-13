package repo

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
	"log/slog"
)

const (
	getUserById         = "SELECT description, username, password_hash, create_time, image_path FROM users WHERE id = $1;"
	updateProfile       = "UPDATE users SET description = $1, password_hash = $2 WHERE id = $3;"
	updateProfileAvatar = "UPDATE users SET image_path = $1 WHERE id = $2;"
)

type ProfileRepo struct {
	db     pgxtype.Querier
	logger *slog.Logger
}

func CreateProfileRepo(db pgxtype.Querier, logger *slog.Logger) *ProfileRepo {
	return &ProfileRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *ProfileRepo) UpdateProfile(ctx context.Context, user models.User) error {
	logger := repo.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	_, err := repo.db.Exec(ctx, updateProfile, user.Description, user.PasswordHash, user.Id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *ProfileRepo) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, imagePath string) error {
	logger := repo.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	_, err := repo.db.Exec(ctx, updateProfileAvatar, imagePath, userID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *ProfileRepo) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
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
