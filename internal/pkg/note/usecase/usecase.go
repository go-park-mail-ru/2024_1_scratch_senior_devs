package usecase

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note"
)

type NoteUsecase struct {
	baseRepo    note.NoteBaseRepo
	searchRepo  note.NoteSearchRepo
	cfg         config.ElasticConfig
	constraints config.ConstraintsConfig
	wg          *sync.WaitGroup
}

func CreateNoteUsecase(baseRepo note.NoteBaseRepo, searchRepo note.NoteSearchRepo, cfg config.ElasticConfig, constraints config.ConstraintsConfig, wg *sync.WaitGroup) *NoteUsecase {
	return &NoteUsecase{
		baseRepo:    baseRepo,
		searchRepo:  searchRepo,
		cfg:         cfg,
		constraints: constraints,
		wg:          wg,
	}
}

func (uc *NoteUsecase) GetAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, searchValue string, tags []string) ([]models.NoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	var res []models.NoteResponse
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

	for i, response := range res {
		info, err := uc.baseRepo.GetOwnerInfo(ctx, response.OwnerId)
		if err != nil {
			logger.Error(err.Error())
		}

		res[i].OwnerInfo = info
	}

	logger.Info("success")
	return res, nil
}

func (uc *NoteUsecase) GetNote(ctx context.Context, noteId uuid.UUID, userId uuid.UUID) (models.NoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteId, userId)
	if err != nil {
		logger.Error(err.Error())
		return models.NoteResponse{}, errors.New("not found")
	}

	if resultNote.OwnerId != userId && !slices.Contains(resultNote.Collaborators, userId) {
		logger.Error("not owner and not collaborator")
		return models.NoteResponse{}, errors.New("not found")
	}

	info, err := uc.baseRepo.GetOwnerInfo(ctx, resultNote.OwnerId)
	if err != nil {
		logger.Error(err.Error())
	}

	resultNote.OwnerInfo = info

	logger.Info("success")
	return resultNote, nil
}

func (uc *NoteUsecase) GetPublicNote(ctx context.Context, noteId uuid.UUID) (models.NoteResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadPublicNote(ctx, noteId)
	if err != nil {
		logger.Error(err.Error())
		return models.NoteResponse{}, errors.New("not found")
	}

	if !resultNote.Public {
		logger.Error("not public")
		return models.NoteResponse{}, errors.New("not found")
	}

	info, err := uc.baseRepo.GetOwnerInfo(ctx, resultNote.OwnerId)
	if err != nil {
		logger.Error(err.Error())
	}

	resultNote.OwnerInfo = info

	logger.Info("success")
	return resultNote, nil
}

func (uc *NoteUsecase) CreateNote(ctx context.Context, userId uuid.UUID, noteData string) (models.Note, error) {
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
		Favorite:      false,
	}

	if err := uc.baseRepo.CreateNote(ctx, newNote); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.CreateNote(ctx, newNote); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return newNote, nil
}

func (uc *NoteUsecase) UpdateNote(ctx context.Context, noteId uuid.UUID, userId uuid.UUID, noteData string) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.baseRepo.ReadNote(ctx, noteId, userId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if updatedNote.OwnerId != userId && !slices.Contains(updatedNote.Collaborators, userId) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not found")
	}

	if updatedNote.Data == noteData {
		logger.Info("note data not modified")
		return updatedNote.Note, nil
	}

	updatedNote.UpdateTime = time.Now().UTC()
	updatedNote.Data = noteData

	if err := uc.baseRepo.UpdateNote(ctx, updatedNote.Note); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.UpdateNote(ctx, updatedNote.Note); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return updatedNote.Note, nil
}

func (uc *NoteUsecase) DeleteNote(ctx context.Context, noteId uuid.UUID, ownerId uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	deletedNote, err := uc.baseRepo.ReadNote(ctx, noteId, ownerId)
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

func (uc *NoteUsecase) getDepth(ctx context.Context, parentParentID uuid.UUID, currentDepth int, userId uuid.UUID) (int, error) {
	emptyID := uuid.UUID{}
	if parentParentID == emptyID {
		return currentDepth, nil
	}

	parent, err := uc.baseRepo.ReadNote(ctx, parentParentID, userId)
	if err != nil {
		return -1, err
	}

	return uc.getDepth(ctx, parent.Parent, currentDepth+1, userId)
}

func (uc *NoteUsecase) CreateSubNote(ctx context.Context, userId uuid.UUID, noteData string, parentID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	newNote := models.Note{
		Id:            uuid.NewV4(),
		Data:          noteData,
		CreateTime:    time.Now().UTC(),
		UpdateTime:    time.Now().UTC(),
		OwnerId:       userId,
		Parent:        parentID,
		Children:      []uuid.UUID{},
		Tags:          []string{},
		Collaborators: []uuid.UUID{},
		Favorite:      false,
	}

	parent, err := uc.baseRepo.ReadNote(ctx, parentID, userId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	if parent.OwnerId != userId && !slices.Contains(parent.Collaborators, userId) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not found")
	}

	if len(parent.Children) >= uc.constraints.MaxSubnotes {
		logger.Error(note.ErrTooManySubnotes)
		return models.Note{}, errors.New(note.ErrTooManySubnotes)
	}

	depth, err := uc.getDepth(ctx, parent.Parent, 1, userId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	if depth >= uc.constraints.MaxDepth {
		logger.Error(note.ErrTooDeep)
		return models.Note{}, errors.New(note.ErrTooDeep)
	}

	newNote.OwnerId = parent.OwnerId
	newNote.Collaborators = parent.Collaborators

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

func (uc *NoteUsecase) addCollaboratorRecursive(ctx context.Context, noteID uuid.UUID, guestID uuid.UUID, userID uuid.UUID) error {
	currentNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		return err
	}

	_, err = uc.baseRepo.AddCollaborator(ctx, noteID, guestID)
	if err != nil {
		return err
	}

	for _, child := range currentNote.Children {
		if err := uc.addCollaboratorRecursive(ctx, child, guestID, userID); err != nil {
			return err
		}
	}

	return nil
}

func (uc *NoteUsecase) AddCollaborator(ctx context.Context, noteID uuid.UUID, userID uuid.UUID, guestID uuid.UUID) (string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	currentNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	if currentNote.OwnerId != userID {
		logger.Error("not owner")
		return "", errors.New("not found")
	}

	emptyID := uuid.UUID{}
	if currentNote.Parent != emptyID {
		logger.Error("note has a parent")
		return "", errors.New("not found")
	}

	if slices.Contains(currentNote.Collaborators, guestID) {
		logger.Error(note.ErrAlreadyCollaborator)
		return "", errors.New(note.ErrAlreadyCollaborator)
	}

	if len(currentNote.Collaborators) >= uc.constraints.MaxCollaborators {
		logger.Error(note.ErrTooManyCollaborators)
		return "", errors.New(note.ErrTooManyCollaborators)
	}

	title, err := uc.baseRepo.AddCollaborator(ctx, noteID, guestID)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	for _, child := range currentNote.Children {
		if err := uc.addCollaboratorRecursive(ctx, child, guestID, userID); err != nil {
			logger.Error(err.Error())
			return "", err
		}
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
	return title, nil
}

func (uc *NoteUsecase) AddTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.baseRepo.ReadNote(ctx, noteId, userId)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	if updatedNote.OwnerId != userId {
		logger.Error("not owner")
		return models.Note{}, errors.New("not found")
	}

	if len(updatedNote.Tags) >= uc.constraints.MaxTags {
		logger.Error(note.ErrTooManyTags)
		return models.Note{}, errors.New(note.ErrTooManyTags)
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

		if err := uc.baseRepo.RememberTag(ctx, tagName, userId); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return updatedNote.Note, nil
}

func (uc *NoteUsecase) DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID, userId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	updatedNote, err := uc.baseRepo.ReadNote(ctx, noteId, userId)
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
	return updatedNote.Note, nil
}

func (uc *NoteUsecase) RememberTag(ctx context.Context, tagName string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := uc.baseRepo.RememberTag(ctx, tagName, userID); err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (uc *NoteUsecase) ForgetTag(ctx context.Context, tagName string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := uc.baseRepo.ForgetTag(ctx, tagName, userID); err != nil {
		logger.Error(err.Error())
		return err
	}

	if err := uc.baseRepo.DeleteTagFromAllNotes(ctx, tagName, userID); err != nil {
		logger.Error(err.Error())
		return err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.DeleteTagFromAllNotes(ctx, tagName, userID); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return nil
}
func (uc *NoteUsecase) UpdateTag(ctx context.Context, oldTag string, newTag string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := uc.baseRepo.UpdateTag(ctx, oldTag, newTag, userID); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := uc.baseRepo.UpdateTagOnAllNotes(ctx, oldTag, newTag, userID); err != nil {
		logger.Error(err.Error())
		return err
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.UpdateTagOnAllNotes(ctx, oldTag, newTag, userID); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return nil
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

func (uc *NoteUsecase) SetIcon(ctx context.Context, noteID uuid.UUID, icon string, userID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not owner and not collaborator")
	}

	if err := uc.baseRepo.SetIcon(ctx, noteID, icon); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	resultNote.Icon = icon

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.SetIcon(ctx, noteID, icon); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return resultNote.Note, nil
}

func (uc *NoteUsecase) SetHeader(ctx context.Context, noteID uuid.UUID, header string, userID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not owner and not collaborator")
	}

	if err := uc.baseRepo.SetHeader(ctx, noteID, header); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	resultNote.Header = header

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.SetHeader(ctx, noteID, header); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return resultNote.Note, nil
}
func (uc *NoteUsecase) AddFav(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not owner and not collaborator")
	}
	if resultNote.Favorite {
		return resultNote.Note, nil
	}
	if err := uc.baseRepo.AddFav(ctx, noteID, userID); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	resultNote.Favorite = true

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.ChangeFlag(ctx, noteID, true); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return resultNote.Note, nil
}

func (uc *NoteUsecase) DelFav(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not owner and not collaborator")
	}
	if !resultNote.Favorite {
		return resultNote.Note, nil
	}
	if err := uc.baseRepo.DelFav(ctx, noteID, userID); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	resultNote.Favorite = false

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.ChangeFlag(ctx, noteID, false); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return resultNote.Note, nil
}

func (uc *NoteUsecase) changeModeRecursive(ctx context.Context, noteID uuid.UUID, userID uuid.UUID, isPublic bool) error {
	currentNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		return err
	}

	if isPublic {
		if err = uc.baseRepo.SetPublic(ctx, noteID); err != nil {
			return err
		}
	} else {
		if err = uc.baseRepo.SetPrivate(ctx, noteID); err != nil {
			return err
		}
	}

	for _, child := range currentNote.Children {
		if err := uc.changeModeRecursive(ctx, child, userID, isPublic); err != nil {
			return err
		}
	}

	return nil
}

func (uc *NoteUsecase) SetPublic(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not owner and not collaborator")
	}

	if err := uc.baseRepo.SetPublic(ctx, noteID); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	resultNote.Public = true

	for _, child := range resultNote.Children {
		if err := uc.changeModeRecursive(ctx, child, userID, true); err != nil {
			logger.Error(err.Error())
			return models.Note{}, err
		}
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.SetPublic(ctx, noteID); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return resultNote.Note, nil
}

func (uc *NoteUsecase) SetPrivate(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return models.Note{}, errors.New("not owner and not collaborator")
	}

	if err := uc.baseRepo.SetPrivate(ctx, noteID); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	resultNote.Public = false

	for _, child := range resultNote.Children {
		if err := uc.changeModeRecursive(ctx, child, userID, false); err != nil {
			logger.Error(err.Error())
			return models.Note{}, err
		}
	}

	uc.wg.Add(1)
	go func() {
		defer uc.wg.Done()
		if err := uc.searchRepo.SetPrivate(ctx, noteID); err != nil {
			logger.Error(err.Error())
		}
	}()
	uc.wg.Wait()

	logger.Info("success")
	return resultNote.Note, nil
}

func (uc *NoteUsecase) GetAttachList(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) ([]string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return []string{}, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return []string{}, errors.New("not owner and not collaborator")
	}

	paths, err := uc.baseRepo.GetAttachList(ctx, noteID)
	if err != nil {
		logger.Error(err.Error())
		return []string{}, err
	}

	logger.Info("success")
	return paths, nil
}

func (uc *NoteUsecase) GetSharedAttachList(ctx context.Context, noteID uuid.UUID) ([]string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadPublicNote(ctx, noteID)
	if err != nil {
		logger.Error(err.Error())
		return []string{}, errors.New("not found")
	}

	if !resultNote.Public {
		logger.Error("not public note")
		return []string{}, errors.New("not found")
	}

	paths, err := uc.baseRepo.GetAttachList(ctx, noteID)
	if err != nil {
		logger.Error(err.Error())
		return []string{}, err
	}

	logger.Info("success")
	return paths, nil
}

func (uc *NoteUsecase) CheckPermissions(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) (bool, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote, err := uc.baseRepo.ReadNote(ctx, noteID, userID)
	if err != nil {
		logger.Error(err.Error())
		return false, errors.New("not found")
	}

	if resultNote.OwnerId != userID && !slices.Contains(resultNote.Collaborators, userID) {
		logger.Error("not owner and not collaborator")
		return false, errors.New("not owner and not collaborator")
	}

	logger.Info("success")
	return true, nil
}
