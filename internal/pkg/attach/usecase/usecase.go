package usecase

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/filework"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"
)

var ErrNoteNotFound = errors.New("note not found")

type AttachUsecase struct {
	repo     attach.AttachRepo
	noteRepo note.NoteBaseRepo
	logger   *slog.Logger
}

func CreateAttachUsecase(repo attach.AttachRepo, noteRepo note.NoteBaseRepo, logger *slog.Logger) *AttachUsecase {
	return &AttachUsecase{
		repo:     repo,
		noteRepo: noteRepo,
		logger:   logger,
	}
}

func (uc *AttachUsecase) AddAttach(ctx context.Context, noteID uuid.UUID, userID uuid.UUID, attach io.ReadSeeker, extension string) (models.Attach, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	noteData, err := uc.noteRepo.ReadNote(ctx, noteID)
	if err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}
	if noteData.OwnerId != userID {
		logger.Error("user is not owner of note")
		return models.Attach{}, ErrNoteNotFound
	}

	attachBasePath := os.Getenv("ATTACHES_BASE_PATH")
	newAttachId := uuid.NewV4()
	newAttachPathNoExtension := newAttachId.String() //+ extension

	newExtension, err := filework.WriteFileOnDisk(path.Join(attachBasePath, newAttachPathNoExtension), extension, attach)
	if err != nil {
		logger.Error("write on disk: " + err.Error())
		return models.Attach{}, err
	}
	newAttachPath := newAttachPathNoExtension + newExtension
	newAttach := models.Attach{
		Id:     newAttachId,
		Path:   newAttachPath,
		NoteId: noteID,
	}

	if err := uc.repo.AddAttach(ctx, newAttach); err != nil {
		logger.Error(err.Error())
		return newAttach, err
	}

	logger.Info("success")
	return newAttach, nil
}

func (uc *AttachUsecase) DeleteAttach(ctx context.Context, attachID uuid.UUID, userID uuid.UUID) error {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	attachData, err := uc.repo.GetAttach(ctx, attachID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	noteData, err := uc.noteRepo.ReadNote(ctx, attachData.NoteId)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if noteData.OwnerId != userID {
		logger.Error("user is not owner of note")
		return ErrNoteNotFound
	}

	err = uc.repo.DeleteAttach(ctx, attachID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if err := os.Remove(path.Join(os.Getenv("ATTACHES_BASE_PATH"), attachData.Path)); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("success")
	return nil
}

func (uc *AttachUsecase) GetAttach(ctx context.Context, attachID uuid.UUID, userID uuid.UUID) (models.Attach, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	attachData, err := uc.repo.GetAttach(ctx, attachID)
	if err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}

	noteData, err := uc.noteRepo.ReadNote(ctx, attachData.NoteId)
	if err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}

	if noteData.OwnerId != userID {
		logger.Error("user is not owner of note")
		return models.Attach{}, ErrNoteNotFound
	}
	return attachData, nil
}
