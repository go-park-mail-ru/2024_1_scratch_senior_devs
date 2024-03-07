package usecase

import (
	"context"
	"github.com/satori/uuid"
	"time"

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

func (uc *NoteUsecase) GetAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, titleSubstr string) ([]models.Note, error) {
	res, err := uc.repo.ReadAllNotes(ctx, userId, count, offset, titleSubstr)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (uc *NoteUsecase) CreateNote(ctx context.Context, userId uuid.UUID) (models.Note, error) {
	newNote := models.Note{
		Id:         uuid.NewV4(),
		Data:       []byte(`{"title":"Новая заметка","content":""}`),
		CreateTime: time.Now().UTC(),
		UpdateTime: nil,
		OwnerId:    userId,
	}

	err := uc.repo.CreateNote(ctx, newNote)
	if err != nil {
		return models.Note{}, err
	}

	return newNote, nil
}
