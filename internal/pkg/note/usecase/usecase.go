package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
)

type NoteUsecase struct {
	repo   note.NoteRepo
	logger *slog.Logger
}

func CreateNoteUsecase(repo note.NoteRepo, logger *slog.Logger) *NoteUsecase {
	return &NoteUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *NoteUsecase) GetAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, titleSubstr string) ([]models.Note, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	res, err := uc.repo.ReadAllNotes(ctx, userId, count, offset, titleSubstr)
	if err != nil {
		logger.Error(err.Error())
		return res, err
	}

	logger.Info("success")
	return res, nil
}

func (uc *NoteUsecase) GetNote(ctx context.Context, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	resultNote, err := uc.repo.ReadNote(ctx, noteId)
	if err != nil || resultNote.OwnerId != userId {
		logger.Error(err.Error())
		return models.Note{}, errors.New("note not found")
	}

	logger.Info("success")
	return resultNote, nil
}

func (uc *NoteUsecase) CreateNote(ctx context.Context, userId uuid.UUID, noteData []byte) (models.Note, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	newNote := models.Note{
		Id:         uuid.NewV4(),
		Data:       noteData,
		CreateTime: time.Now().UTC(),
		UpdateTime: time.Now().UTC(),
		OwnerId:    userId,
	}

	if err := uc.repo.CreateNote(ctx, newNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return newNote, nil
}

func (uc *NoteUsecase) UpdateNote(ctx context.Context, id uuid.UUID, noteData []byte) (models.Note, error) {

	updatedNote, err := uc.repo.ReadNote(ctx, id)
	if err != nil {
		return models.Note{}, err
	}

	updatedNote.UpdateTime = time.Now().UTC()
	updatedNote.Data = noteData //data.Data

	err = uc.repo.UpdateNote(ctx, updatedNote)
	if err != nil {
		return models.Note{}, err
	}

	return updatedNote, nil

}
func (uc *NoteUsecase) DeleteNote(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	err := uc.repo.DeleteNote(ctx, id, ownerId)
	if err != nil {
		return err
	}
	return nil

}
