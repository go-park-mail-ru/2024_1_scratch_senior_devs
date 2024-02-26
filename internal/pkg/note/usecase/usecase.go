package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/repo"
	"github.com/satori/uuid"
)

type NotesUsecase struct {
	repo repo.NotesRepo
}

func CreateNotesUsecase(repo repo.NotesRepo) *NotesUsecase {
	return &NotesUsecase{
		repo: repo,
	}
}

func (uc *NotesUsecase) GetAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64) ([]models.Note, error) {
	res, err := uc.repo.ReadAllNotes(ctx, userId, count, offset)
	if err != nil {
		return nil, err
	}
	return res, nil
}
