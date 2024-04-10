package usecase

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/protection"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/code"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/sources"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth"
)

type AuthUsecase struct {
	repo          auth.AuthRepo
	noteRepo      note.NoteRepo
	logger        *slog.Logger
	cfg           config.AuthUsecaseConfig
	cfgValidation config.UserValidationConfig
}

func CreateAuthUsecase(repo auth.AuthRepo, noteRepo note.NoteRepo, logger *slog.Logger, cfg config.AuthUsecaseConfig, cfgValidation config.UserValidationConfig) *AuthUsecase {
	return &AuthUsecase{
		repo:          repo,
		noteRepo:      noteRepo,
		logger:        logger,
		cfg:           cfg,
		cfgValidation: cfgValidation,
	}
}

func (uc *AuthUsecase) SignUp(ctx context.Context, data models.UserFormData) (models.User, string, time.Time, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(uc.cfg.JWTLifeTime)

	newUser := models.User{
		Id:           uuid.NewV4(),
		Username:     data.Username,
		PasswordHash: responses.GetHash(data.Password),
		ImagePath:    uc.cfg.DefaultImagePath,
		CreateTime:   currentTime,
		SecondFactor: "",
	}

	if err := uc.repo.CreateUser(ctx, newUser); err != nil {
		logger.Error(err.Error())
		return models.User{}, "", currentTime, auth.ErrCreatingUser
	}

	jwtToken, err := protection.GenJwtToken(newUser, uc.cfg.JWTLifeTime)
	if err != nil {
		logger.Error("GenJwtToken error: " + err.Error())
		return models.User{}, "", currentTime, err
	}

	if err := uc.noteRepo.CreateNote(ctx, models.Note{
		Id:         uuid.NewV4(),
		Data:       uc.noteRepo.MakeHelloNoteData(newUser.Username),
		CreateTime: currentTime,
		UpdateTime: currentTime,
		OwnerId:    newUser.Id,
	}); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("success")
	return newUser, jwtToken, expTime, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, data models.UserFormData) (models.User, string, time.Time, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	currentTime := time.Now().UTC()
	expTime := currentTime.Add(uc.cfg.JWTLifeTime)

	user, err := uc.repo.GetUserByUsername(ctx, data.Username)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, "", currentTime, auth.ErrUserNotFound
	}
	if user.PasswordHash != responses.GetHash(data.Password) {
		logger.Error("wrong password")
		return models.User{}, "", currentTime, auth.ErrWrongUserData
	}

	if user.SecondFactor != "" {
		if data.Code == "" {
			logger.Error(auth.ErrFirstFactorPassed.Error())
			return models.User{}, "", currentTime, auth.ErrFirstFactorPassed
		}

		err := code.CheckCode(data.Code, string(user.SecondFactor))
		if err != nil {
			logger.Error(err.Error())
			return models.User{}, "", currentTime, auth.ErrWrongAuthCode
		}
	}

	jwtToken, err := protection.GenJwtToken(user, uc.cfg.JWTLifeTime)
	if err != nil {
		logger.Error("GenJwtToken error: " + err.Error())
		return models.User{}, "", currentTime, err
	}

	logger.Info("success")
	return user, jwtToken, expTime, nil
}

func (uc *AuthUsecase) CheckUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	userData, err := uc.repo.GetUserById(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, auth.ErrUserNotFound
	}

	logger.Info("success")
	return userData, nil
}

func (uc *AuthUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, payload models.ProfileUpdatePayload) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	payload.Sanitize()

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, auth.ErrUserNotFound
	}

	if payload.Password.Old != "" && payload.Password.New != "" {
		if err := payload.Validate(uc.cfgValidation); err != nil {
			logger.Error("validation error: " + err.Error())
			return models.User{}, err
		}

		if user.PasswordHash != responses.GetHash(payload.Password.Old) {
			logger.Error("wrong password")
			return models.User{}, auth.ErrWrongPassword
		}

		user.PasswordHash = responses.GetHash(payload.Password.New)
	}

	user.Description = payload.Description

	if err := uc.repo.UpdateProfile(ctx, user); err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	logger.Info("success")
	return user, nil
}

func (uc *AuthUsecase) UpdateProfileAvatar(ctx context.Context, userID uuid.UUID, avatar io.ReadSeeker, extension string) (models.User, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.User{}, auth.ErrUserNotFound
	}

	imagesBasePath := os.Getenv("IMAGES_BASE_PATH")
	imagePathNoExtension := uuid.NewV4().String() //+ ".webp" // + extension

	newExtension, err := sources.WriteFileOnDisk(path.Join(imagesBasePath, imagePathNoExtension), extension, avatar)
	if err != nil {
		logger.Error("write on disk: " + err.Error())
		return models.User{}, err
	}
	newImagePath := imagePathNoExtension + newExtension
	if err := uc.repo.UpdateProfileAvatar(ctx, userID, newImagePath); err != nil {
		logger.Error(err.Error())
		return models.User{}, err
	}

	// удаление старой аватарки делаем только после успешного создания новой
	if user.ImagePath != "default.jpg" {
		if err := os.Remove(path.Join(imagesBasePath, user.ImagePath)); err != nil {
			logger.Error(err.Error())
		}
	}

	user.ImagePath = newImagePath

	logger.Info("success")
	return user, nil
}

func (uc *AuthUsecase) GenerateAndUpdateSecret(ctx context.Context, username string) ([]byte, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	secret := code.GenerateSecret()
	if err := uc.repo.UpdateSecret(ctx, username, string(secret)); err != nil {
		logger.Error(err.Error())
		return []byte{}, err
	}

	logger.Info("success")
	return secret, nil
}

func (uc *AuthUsecase) DeleteSecret(ctx context.Context, username string) error {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	if err := uc.repo.DeleteSecret(ctx, username); err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
