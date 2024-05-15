package repo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

const (
	getAllNotes = `
	SELECT id, data, create_time, update_time, owner_id, parent, children, tags, collaborators, icon, header, favorite FROM notes
	WHERE parent = '00000000-0000-0000-0000-000000000000'::UUID
	AND (
		owner_id = $1
		OR $1 = ANY(collaborators)
	)
	AND (
		cardinality($4::TEXT[]) = 0 OR $4::TEXT[] IS NULL OR array(select unnest($4::TEXT[]) except select unnest(tags)) = '{}'
	)
	ORDER BY favorite DESC, update_time DESC
	LIMIT $2 OFFSET $3;
	`
	getNote    = `SELECT id, data, create_time, update_time, owner_id, parent, children, tags, collaborators, icon, header, favorite FROM notes WHERE id = $1;`
	createNote = "INSERT INTO notes(id, data, create_time, update_time, owner_id, parent, children, tags, collaborators, icon, header, favorite) VALUES ($1, $2::json, $3, $4, $5, $6, $7::UUID[], $8::TEXT[], $9::UUID[], $10, $11, $12);"
	updateNote = "UPDATE notes SET data = $1, update_time = $2 WHERE id = $3; "
	deleteNote = "DELETE FROM notes CASCADE WHERE id = $1;"

	addSubNote    = "UPDATE notes SET children = array_append(children, $1) WHERE id = $2;"
	removeSubNote = "UPDATE notes SET children = array_remove(children, $1) WHERE id = $2;"

	getUpdates = "SELECT note_id, created, message_info FROM messages WHERE note_id = $1 AND created > $2;"

	addCollaborator = "UPDATE notes SET collaborators = array_append(collaborators, $1) WHERE id = $2;"

	addTag    = "UPDATE notes SET tags = array_append(tags, $1) WHERE id = $2;"
	deleteTag = "UPDATE notes SET tags = array_remove(tags, $1) WHERE id = $2;"

	getTags               = `SELECT tag_name FROM all_tags WHERE user_id = $1;`
	rememberTag           = "INSERT INTO all_tags(tag_name, user_id) VALUES ($1, $2) ON CONFLICT (tag_name, user_id) DO NOTHING;"
	forgetTag             = "DELETE FROM all_tags WHERE tag_name = $1 AND user_id = $2;"
	updateTag             = "UPDATE all_tags SET tag_name = $1 WHERE user_id = $2 AND tag_name = $3;"
	deleteTagFromAllNotes = "UPDATE notes SET tags = array_remove(tags, $1) WHERE owner_id = $2;"

	setIcon   = "UPDATE notes SET icon = $1 WHERE id = $2;"
	setHeader = "UPDATE notes SET header = $1 WHERE id = $2;"

	changeFlag = "UPDATE notes SET favorite = $1 WHERE id = $2;"
)

type NotePostgres struct {
	db   pgxtype.Querier
	metr metrics.DBMetrics
}

func CreateNotePostgres(db pgxtype.Querier, metr metrics.DBMetrics) *NotePostgres {
	return &NotePostgres{
		db:   db,
		metr: metr,
	}
}

func (repo *NotePostgres) GetTags(ctx context.Context, userID uuid.UUID) ([]string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result := make([]string, 0)

	start := time.Now()
	query, err := repo.db.Query(ctx, getTags, userID)
	repo.metr.ObserveResponseTime("getTags", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("getTags")
		return result, err
	}

	for query.Next() {
		var tag string
		if err := query.Scan(&tag); err != nil {
			logger.Error(err.Error())
			return result, fmt.Errorf("error occured while scanning tags: %w", err)
		}
		result = append(result, tag)
	}

	return result, nil
}

func (repo *NotePostgres) AddCollaborator(ctx context.Context, noteID uuid.UUID, guestID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, addCollaborator, guestID, noteID)
	repo.metr.ObserveResponseTime("addCollaborator", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("addCollaborator")
		return err
	}

	return nil
}

func (repo *NotePostgres) GetUpdates(ctx context.Context, noteID uuid.UUID, offset time.Time) ([]models.Message, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result := make([]models.Message, 0)

	start := time.Now()
	query, err := repo.db.Query(ctx, getUpdates, noteID, offset)
	repo.metr.ObserveResponseTime("getUpdates", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("getUpdates")
		return result, err
	}

	for query.Next() {
		var message models.Message
		if err := query.Scan(&message.NoteId, &message.Created, &message.MessageInfo); err != nil {
			logger.Error(err.Error())
			return result, fmt.Errorf("error occured while scanning messages: %w", err)
		}
		result = append(result, message)
	}

	return result, nil
}

func (repo *NotePostgres) ReadAllNotes(ctx context.Context, userId uuid.UUID, count int64, offset int64, tags []string) ([]models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	result := make([]models.Note, 0, count)

	start := time.Now()
	query, err := repo.db.Query(ctx, getAllNotes, userId, count, offset, tags)
	repo.metr.ObserveResponseTime("getAllNotes", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("getAllNotes")
		return result, err
	}

	for query.Next() {
		var note models.Note
		if err := query.Scan(
			&note.Id,
			&note.Data,
			&note.CreateTime,
			&note.UpdateTime,
			&note.OwnerId,
			&note.Parent,
			&note.Children,
			&note.Tags,
			&note.Collaborators,
			&note.Icon,
			&note.Header,
			&note.Favorite,
		); err != nil {
			logger.Error("scanning" + err.Error())
			return result, fmt.Errorf("error occured while scanning notes: %w", err)
		}
		result = append(result, note)
	}

	logger.Info("success")
	return result, nil
}

func (repo *NotePostgres) ReadNote(ctx context.Context, noteId uuid.UUID) (models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	resultNote := models.Note{}

	start := time.Now()
	err := repo.db.QueryRow(ctx, getNote, noteId).Scan(
		&resultNote.Id,
		&resultNote.Data,
		&resultNote.CreateTime,
		&resultNote.UpdateTime,
		&resultNote.OwnerId,
		&resultNote.Parent,
		&resultNote.Children,
		&resultNote.Tags,
		&resultNote.Collaborators,
		&resultNote.Icon,
		&resultNote.Header,
		&resultNote.Favorite,
	)
	repo.metr.ObserveResponseTime("getNote", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("getNote")
		return models.Note{}, err
	}

	logger.Info("success")
	return resultNote, nil
}

func (repo *NotePostgres) CreateNote(ctx context.Context, note models.Note) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, createNote,
		note.Id,
		note.Data,
		note.CreateTime,
		note.UpdateTime,
		note.OwnerId,
		note.Parent,
		note.Children,
		note.Tags,
		note.Collaborators,
		note.Icon,
		note.Header,
		note.Favorite,
	)
	repo.metr.ObserveResponseTime("createNote", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("createNote")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) UpdateNote(ctx context.Context, note models.Note) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, updateNote, note.Data, note.UpdateTime, note.Id)
	repo.metr.ObserveResponseTime("updateNote", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("updateNote")
		return err
	}

	logger.Info("success")
	return nil

}

func (repo *NotePostgres) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, deleteNote, id)
	repo.metr.ObserveResponseTime("deleteNote", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("deleteNote")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) AddSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, addSubNote, childID, id)
	repo.metr.ObserveResponseTime("addSubNote", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("addSubNote")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) RemoveSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, removeSubNote, childID, id)
	repo.metr.ObserveResponseTime("removeSubNote", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("removeSubNote")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) AddTag(ctx context.Context, tagName string, noteId uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, addTag, tagName, noteId)
	repo.metr.ObserveResponseTime("addTag", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("addTag")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, deleteTag, tagName, noteId)
	repo.metr.ObserveResponseTime("deleteTag", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("deleteTag")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) RememberTag(ctx context.Context, tagName string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	result, err := repo.db.Exec(ctx, rememberTag, tagName, userID)
	repo.metr.ObserveResponseTime("rememberTag", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("rememberTag")
		return err
	}

	if result.RowsAffected() == 0 {
		logger.Error("tag already exists")
		return errors.New("tag already exists")
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) ForgetTag(ctx context.Context, tagName string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, forgetTag, tagName, userID)
	repo.metr.ObserveResponseTime("forgetTag", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("forgetTag")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) DeleteTagFromAllNotes(ctx context.Context, tagName string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, deleteTagFromAllNotes, tagName, userID)
	repo.metr.ObserveResponseTime("deleteTagFromAllNotes", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("deleteTagFromAllNotes")
		return err
	}

	logger.Info("success")
	return nil
}
func (repo *NotePostgres) UpdateTag(ctx context.Context, oldTag string, newTag string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, updateTag, newTag, userID, oldTag)
	repo.metr.ObserveResponseTime("updateTag", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("updateTag")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) SetIcon(ctx context.Context, noteID uuid.UUID, icon string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, setIcon, icon, noteID)
	repo.metr.ObserveResponseTime("setIcon", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("setIcon")
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NotePostgres) SetHeader(ctx context.Context, noteID uuid.UUID, header string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, setHeader, header, noteID)
	repo.metr.ObserveResponseTime("setHeader", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("setHeader")
		return err
	}

	logger.Info("success")
	return nil

}

func (repo *NotePostgres) ChangeFlag(ctx context.Context, noteID uuid.UUID, flag bool) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.db.Exec(ctx, changeFlag, flag, noteID)
	repo.metr.ObserveResponseTime("changeFlag", time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors("changeFlag")
		return err
	}

	logger.Info("success")
	return nil
}
