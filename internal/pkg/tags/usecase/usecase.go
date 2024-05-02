package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/tags"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"
)

type TagUsecase struct {
	repo     tags.TagRepo
	noteRepo note.NoteBaseRepo
}

func CreateTagUsecase(repo tags.TagRepo, noteRepo note.NoteBaseRepo) *TagUsecase {
	return &TagUsecase{
		repo:     repo,
		noteRepo: noteRepo,
	}
}

func (uc *TagUsecase) AddTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.noteRepo.ReadNote(ctx, noteId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	if updatedNote.OwnerId != userId {

		logger.Error("not owner")
		return models.Note{}, errors.New("not found")

	}

	if err := uc.repo.AddTag(ctx, tagName, noteId); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	updatedNote.Tags = append(updatedNote.Tags, tagName)
	logger.Info("success")
	return updatedNote, nil
}

func (uc *TagUsecase) DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.noteRepo.ReadNote(ctx, noteId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	if updatedNote.OwnerId != userId {

		logger.Error("not owner")
		return models.Note{}, errors.New("not found")

	}

	if err := uc.repo.DeleteTag(ctx, tagName, noteId); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	newTags := make([]string, 0)
	for i := range updatedNote.Tags {
		if updatedNote.Tags[i] != tagName {
			newTags = append(newTags, updatedNote.Tags[i])
		}
	}
	updatedNote.Tags = newTags
	return updatedNote, nil
}
