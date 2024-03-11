package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils"
	"github.com/satori/uuid"
	"log/slog"
	"time"

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
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	res, err := uc.repo.ReadAllNotes(ctx, userId, count, offset, titleSubstr)
	if err != nil {
		logger.Error(err.Error())
		return res, err
	}

	logger.Info("success")
	return res, nil
}

func (uc *NoteUsecase) GetNote(ctx context.Context, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	resultNote, err := uc.repo.ReadNote(ctx, noteId)
	if err != nil || resultNote.OwnerId != userId {
		logger.Error(err.Error())
		return models.Note{}, errors.New("note not found")
	}

	logger.Info("success")
	return resultNote, nil
}

func (uc *NoteUsecase) CreateNote(ctx context.Context, userId uuid.UUID) (models.Note, error) {
	logger := uc.logger.With(slog.String("ID", utils.GetRequestId(ctx)), slog.String("func", utils.GFN()))

	newNote := models.Note{
		Id:         uuid.NewV4(),
		Data:       []byte(`{"title":"Новая заметка","content":""}`),
		CreateTime: time.Now().UTC(),
		UpdateTime: nil,
		OwnerId:    userId,
	}

	if err := uc.repo.CreateNote(ctx, newNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return newNote, nil
}
