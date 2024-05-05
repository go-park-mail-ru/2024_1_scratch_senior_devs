package usecase

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/elasticsearch"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
)

type NoteUsecase struct {
	baseRepo   note.NoteBaseRepo
	searchRepo note.NoteSearchRepo
	cfg        config.ElasticConfig
	wg         *sync.WaitGroup
}

func CreateNoteUsecase(baseRepo note.NoteBaseRepo, searchRepo note.NoteSearchRepo, cfg config.ElasticConfig, wg *sync.WaitGroup) *NoteUsecase {
	return &NoteUsecase{
		baseRepo:   baseRepo,
		searchRepo: searchRepo,
		cfg:        cfg,
		wg:         wg,
	}
}

func (uc *NoteUsecase) GetAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, searchValue string, tags []string) ([]models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	var res []models.Note
	var err error

	if utf8.RuneCountInString(searchValue) < uc.cfg.ElasticSearchValueMinLength {
		res, err = uc.baseRepo.ReadAllNotes(ctx, userId, count, offset, tags)
	} else {
		res, err = uc.searchRepo.SearchNotes(ctx, userId, count, offset, searchValue, tags)
	}

	if err != nil {
		logger.Error(err.Error())
		return res, err
	}

	logger.Info("success")
	return res, nil
}

func (uc *NoteUsecase) GetNote(ctx context.Context, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if resultNote.OwnerId != userId && !slices.Contains(resultNote.Collaborators, userId) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not found")
	}

	logger.Info("success")
	return resultNote, nil
}

func (uc *NoteUsecase) CreateNote(ctx context.Context, userId uuid.UUID, noteData []byte) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	newNote := models.Note{
		Id:            uuid.NewV4(),
		Data:          noteData,
		CreateTime:    time.Now().UTC(),
		UpdateTime:    time.Now().UTC(),
		OwnerId:       userId,
		Parent:        uuid.UUID{},
		Children:      []uuid.UUID{},
		Tags:          []string{},
		Collaborators: []uuid.UUID{},
	}

	if err := uc.baseRepo.CreateNote(ctx, newNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.CreateNote(ctx, elasticsearch.ConvertToElasticNote(newNote)); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return newNote, nil
}

func (uc *NoteUsecase) UpdateNote(ctx context.Context, noteId uuid.UUID, userId uuid.UUID, noteData []byte) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if updatedNote.OwnerId != userId && !slices.Contains(updatedNote.Collaborators, userId) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not found")
	}

	if bytes.Equal(updatedNote.Data, noteData) {
		logger.Info("note data not modified")
		return updatedNote, nil
	}

	updatedNote.UpdateTime = time.Now().UTC()
	updatedNote.Data = noteData

	if err := uc.baseRepo.UpdateNote(ctx, updatedNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.UpdateNote(ctx, elasticsearch.ConvertToElasticNote(updatedNote)); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return updatedNote, nil
}

func (uc *NoteUsecase) DeleteNote(ctx context.Context, noteId uuid.UUID, ownerId uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	deletedNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil || deletedNote.OwnerId != ownerId {
		logger.Error(err.Error())
		return err
	}

	if err := uc.baseRepo.DeleteNote(ctx, noteId); err != nil {
		logger.Error(err.Error())
		return err
	}

	emptyID := uuid.UUID{}
	if deletedNote.Parent != emptyID {
		if err := uc.baseRepo.RemoveSubNote(ctx, deletedNote.Parent, noteId); err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()

		if err := uc.searchRepo.DeleteNote(ctx, noteId); err != nil {
			logger.Error(err.Error())
		}

		if deletedNote.Parent != emptyID {
			if err := uc.searchRepo.RemoveSubNote(ctx, deletedNote.Parent, noteId); err != nil {
				logger.Error(err.Error())
			}
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return nil
}

func (uc *NoteUsecase) CreateSubNote(ctx context.Context, userId uuid.UUID, noteData []byte, parentID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	newNote := models.Note{
		Id:         uuid.NewV4(),
		Data:       noteData,
		CreateTime: time.Now().UTC(),
		UpdateTime: time.Now().UTC(),
		OwnerId:    userId,
		Parent:     parentID,
		Children:   []uuid.UUID{},
	}

	if err := uc.baseRepo.CreateNote(ctx, newNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	if err := uc.baseRepo.AddSubNote(ctx, parentID, newNote.Id); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.AddSubNote(ctx, parentID, newNote.Id); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return newNote, nil
}

func (uc *NoteUsecase) AddCollaborator(ctx context.Context, noteID uuid.UUID, userID uuid.UUID, guestID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	currentNote, err := uc.baseRepo.ReadNote(ctx, noteID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if currentNote.OwnerId != userID {
		logger.Error("not owner")
		return errors.New("not found")
	}

	emptyID := uuid.UUID{}
	if currentNote.Parent != emptyID {
		logger.Error("note has a parent")
		return errors.New("not found")
	}

	if err := uc.baseRepo.AddCollaborator(ctx, noteID, guestID); err != nil {
		logger.Error(err.Error())
		return err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.AddCollaborator(ctx, noteID, guestID); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return nil
}

func (uc *NoteUsecase) AddTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	if updatedNote.OwnerId != userId {
		logger.Error("not owner")
		return models.Note{}, errors.New("not found")
	}

	if err := uc.baseRepo.AddTag(ctx, tagName, noteId); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	updatedNote.Tags = append(updatedNote.Tags, tagName)

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.AddTag(ctx, tagName, noteId); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return updatedNote, nil
}

func (uc *NoteUsecase) DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.baseRepo.ReadNote(ctx, noteId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	if updatedNote.OwnerId != userId {
		logger.Error("not owner")
		return models.Note{}, errors.New("not found")
	}

	if err := uc.baseRepo.DeleteTag(ctx, tagName, noteId); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	newTags := make([]string, 0)
	for i := range updatedNote.Tags {
		if updatedNote.Tags[i] != tagName {
			newTags = append(newTags, updatedNote.Tags[i])
		}
	}
	updatedNote.Tags = newTags

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.DeleteTag(ctx, tagName, noteId); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return updatedNote, nil
}

func (uc *NoteUsecase) GetTags(ctx context.Context, userID uuid.UUID) ([]string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	tags, err := uc.baseRepo.GetTags(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return []string{}, err
	}

	logger.Info("success")
	return tags, nil
}

func (uc *NoteUsecase) CheckPermissions(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (bool, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID)
	if err != nil {
		logger.Error(err.Error())
		return false, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return false, nil
	}

	logger.Info("success")
	return true, nil
}
