package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"log/slog"
	"time"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/validation"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
)

type NoteUsecase struct {
	baseRepo   note.NoteBaseRepo
	searchRepo note.NoteSearchRepo
	logger     *slog.Logger
	cfg        config.ElasticConfig
}

func CreateNoteUsecase(baseRepo note.NoteBaseRepo, searchRepo note.NoteSearchRepo, logger *slog.Logger, cfg config.ElasticConfig) *NoteUsecase {
	return &NoteUsecase{
		baseRepo:   baseRepo,
		searchRepo: searchRepo,
		logger:     logger,
		cfg:        cfg,
	}
}

func (uc *NoteUsecase) GetAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, searchValue string) ([]models.Note, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	res := make([]models.Note, 0)
	var err error

	if utf8.RuneCountInString(searchValue) < uc.cfg.ElasticSearchValueMinLength {
		res, err = uc.baseRepo.ReadAllNotes(ctx, userId, count, offset)
	} else {
		res, err = uc.searchRepo.SearchNotes(ctx, userId, count, offset, searchValue)
	}

	if err != nil {
		logger.Error(err.Error())
		return res, err
	}

	logger.Info("success")
	return res, nil
}

func (uc *NoteUsecase) GetNote(ctx context.Context, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil || resultNote.OwnerId != userId {
		logger.Error(err.Error())
		return models.Note{}, errors.New("note not found")
	}

	logger.Info("success")
	return resultNote, nil
}

func (uc *NoteUsecase) CreateNote(ctx context.Context, userId uuid.UUID, noteData []byte) (models.Note, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	if err := validation.CheckNoteTitle(noteData); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	newNote := models.Note{
		Id:         uuid.NewV4(),
		Data:       noteData,
		CreateTime: time.Now().UTC(),
		UpdateTime: time.Now().UTC(),
		OwnerId:    userId,
	}

	if err := uc.baseRepo.CreateNote(ctx, newNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	if err := uc.searchRepo.CreateNote(ctx, newNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return newNote, nil
}

func (uc *NoteUsecase) UpdateNote(ctx context.Context, noteId uuid.UUID, ownerId uuid.UUID, noteData []byte) (models.Note, error) {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	if err := validation.CheckNoteTitle(noteData); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	updatedNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil || updatedNote.OwnerId != ownerId {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	updatedNote.UpdateTime = time.Now().UTC()
	updatedNote.Data = noteData

	if err := uc.baseRepo.UpdateNote(ctx, updatedNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	if err := uc.searchRepo.UpdateNote(ctx, updatedNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return updatedNote, nil
}

func (uc *NoteUsecase) DeleteNote(ctx context.Context, noteId uuid.UUID, ownerId uuid.UUID) error {
	logger := uc.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	deletedNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil || deletedNote.OwnerId != ownerId {
		logger.Error(err.Error())
		return err
	}

	if err := uc.baseRepo.DeleteNote(ctx, noteId); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := uc.searchRepo.DeleteNote(ctx, noteId); err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
