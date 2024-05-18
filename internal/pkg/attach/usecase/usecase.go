package usecase

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path"
	"slices"

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
}

func CreateAttachUsecase(repo attach.AttachRepo, noteRepo note.NoteBaseRepo) *AttachUsecase {
	return &AttachUsecase{
		repo:     repo,
		noteRepo: noteRepo,
	}
}

func (uc *AttachUsecase) AddAttach(ctx context.Context, noteID uuid.UUID, userID uuid.UUID, attach io.ReadSeeker, extension string) (models.Attach, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.noteRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}
	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Attach{}, errors.New("not found")
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
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	attachData, err := uc.repo.GetAttach(ctx, attachID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	resultNote, err := uc.noteRepo.ReadNote(ctx, attachData.NoteId, userID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return nil
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
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	attachData, err := uc.repo.GetAttach(ctx, attachID)
	if err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}

	resultNote, err := uc.noteRepo.ReadNote(ctx, attachData.NoteId, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Attach{}, err
	}
	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Attach{}, errors.New("not found")
	}

	return attachData, nil
}
