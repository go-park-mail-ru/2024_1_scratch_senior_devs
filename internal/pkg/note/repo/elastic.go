package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/elasticsearch"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/olivere/elastic/v7"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
)

var (
	ErrCantGetResponse = errors.New("can`t get response")
)

type NoteElastic struct {
	elastic *elastic.Client
	cfg     config.ElasticConfig
}

func CreateNoteElastic(elastic *elastic.Client, cfg config.ElasticConfig) *NoteElastic {
	return &NoteElastic{
		elastic: elastic,
		cfg:     cfg,
	}
}

func (repo *NoteElastic) updateNullFieldToEmptyArray(ctx context.Context, fieldName string, noteID uuid.UUID) error {
	script := elastic.NewScript(fmt.Sprintf("if (ctx._source.%s == null) {ctx._source.%s = []; }", fieldName, fieldName))

	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (repo *NoteElastic) SearchNotes(ctx context.Context, userID uuid.UUID, count int64, offset int64, searchValue string, tags []string) ([]models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	ownerQuery := elastic.NewTermsQuery("owner_id", strings.ToLower(userID.String()))
	collaboratorQuery := elastic.NewTermsQuery("collaborators", strings.ToLower(userID.String()))
	searchQuery := elastic.NewMatchQuery("data", searchValue)

	userIdQuery := elastic.NewBoolQuery().Should(ownerQuery, collaboratorQuery)
	fullQuery := elastic.NewBoolQuery().Must(searchQuery, userIdQuery)

	if len(tags) > 0 {
		tagQueries := make([]elastic.Query, len(tags))
		for i, tag := range tags {
			tagQueries[i] = elastic.NewTermsQuery("tags", tag)
		}

		tagsQuery := elastic.NewBoolQuery().Should(tagQueries...)
		fullQuery = fullQuery.Filter(tagsQuery)
	}

	search, err := repo.elastic.Search().
		Query(fullQuery).
		Index(repo.cfg.ElasticIndexName).
		From(int(offset)).
		Size(int(count)).
		Sort("update_time", false).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return []models.Note{}, ErrCantGetResponse
	}

	notes := make([]models.Note, 0)
	for _, hit := range search.Hits.Hits {
		note := models.ElasticNote{}
		if err := json.Unmarshal(hit.Source, &note); err != nil {
			logger.Error(err.Error())
			return []models.Note{}, err
		}
		notes = append(notes, elasticsearch.ConvertToUsualNote(note))
	}

	logger.Info("success")
	return notes, nil
}

func (repo *NoteElastic) ReadNote(ctx context.Context, noteID uuid.UUID) (models.ElasticNote, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	search, err := repo.elastic.Search().
		Index(repo.cfg.ElasticIndexName).
		Query(elastic.NewTermQuery("_id", noteID)).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return models.ElasticNote{}, ErrCantGetResponse
	}

	if len(search.Hits.Hits) == 0 {
		logger.Error("note not found")
		return models.ElasticNote{}, errors.New("note not found")
	}

	note := models.ElasticNote{}
	if err := json.Unmarshal(search.Hits.Hits[0].Source, &note); err != nil {
		logger.Error(err.Error())
		return models.ElasticNote{}, err
	}

	logger.Info("success")
	return note, nil
}

func (repo *NoteElastic) CreateNote(ctx context.Context, note models.ElasticNote) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	noteJSON, err := json.Marshal(note)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var noteMap map[string]interface{}
	if err := json.Unmarshal(noteJSON, &noteMap); err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = repo.elastic.Index().
		Index(repo.cfg.ElasticIndexName).
		Id(note.Id.String()).
		BodyJson(noteMap).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return ErrCantGetResponse
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) UpdateNote(ctx context.Context, note models.ElasticNote) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	noteJSON, err := json.Marshal(note)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var noteMap map[string]interface{}
	if err := json.Unmarshal(noteJSON, &noteMap); err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(note.Id.String()).
		Doc(noteMap).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return ErrCantGetResponse
	}

	logger.Info("success")
	return nil

}

func (repo *NoteElastic) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	_, err := repo.elastic.Delete().
		Index(repo.cfg.ElasticIndexName).
		Id(id.String()).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return ErrCantGetResponse
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) AddSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.children.add(params.childID)").Lang("painless").Param("childID", childID.String())

	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(id.String()).
		Script(script).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) RemoveSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.children.removeIfContains(params.childID)").Lang("painless").Param("childID", childID.String())

	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(id.String()).
		Script(script).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) AddCollaborator(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.collaborators.add(params.collaboratorID)").Lang("painless").Param("collaboratorID", userID.String())

	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) AddTag(ctx context.Context, tagName string, noteID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := repo.updateNullFieldToEmptyArray(ctx, "tags", noteID); err != nil {
		logger.Error(err.Error())
		return err
	}

	script := elastic.NewScript("ctx._source.tags.add(params.tagName)").Lang("painless").Param("tagName", tagName)

	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) DeleteTag(ctx context.Context, tagName string, noteID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	if err := repo.updateNullFieldToEmptyArray(ctx, "tags", noteID); err != nil {
		logger.Error(err.Error())
		return err
	}

	script := elastic.NewScript("ctx._source.tags.removeIfContains(params.tagName)").Lang("painless").Param("tagName", tagName)

	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
