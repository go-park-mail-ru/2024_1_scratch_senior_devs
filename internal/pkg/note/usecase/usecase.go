package usecase

import (
	"context"

	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
)

type NoteUsecase struct {
	repo note.NoteRepo
}

func CreateNoteUsecase(repo note.NoteRepo) *NoteUsecase {
	return &NoteUsecase{
		repo: repo,
	}
}

func (uc *NoteUsecase) GetAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64) ([]models.Note, error) {
	res, err := uc.repo.ReadAllNotes(ctx, userId, count, offset)
	if err != nil {
		return nil, err
	}
	return res, nil
}
