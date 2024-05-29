package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"log/slog"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
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
	metr    metrics.DBMetrics
}

func CreateNoteElastic(elastic *elastic.Client, cfg config.ElasticConfig, metr metrics.DBMetrics) *NoteElastic {
	return &NoteElastic{
		elastic: elastic,
		cfg:     cfg,
		metr:    metr,
	}
}

func (repo *NoteElastic) updateNullFieldToEmptyArray(ctx context.Context, fieldName string, noteID uuid.UUID) error {
	script := elastic.NewScript(fmt.Sprintf("if (ctx._source.%s == null) {ctx._source.%s = []; }", fieldName, fieldName))

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	return nil
}

func (repo *NoteElastic) SearchNotes(ctx context.Context, userID uuid.UUID, count int64, offset int64, searchValue string, tags []string) ([]models.NoteResponse, error) {
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

	start := time.Now()
	search, err := repo.elastic.Search().
		Query(fullQuery).
		Index(repo.cfg.ElasticIndexName).
		From(int(offset)).
		Size(int(count)).
		Sort("update_time", false).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return []models.NoteResponse{}, ErrCantGetResponse
	}

	notes := make([]models.NoteResponse, 0)
	for _, hit := range search.Hits.Hits {
		note := models.NoteResponse{}
		if err := json.Unmarshal(hit.Source, &note); err != nil {
			logger.Error(err.Error())
			return []models.NoteResponse{}, err
		}
		notes = append(notes, note)
	}

	logger.Info("success")
	return notes, nil
}

func (repo *NoteElastic) CreateNote(ctx context.Context, note models.Note) error {
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

	start := time.Now()
	_, err = repo.elastic.Index().
		Index(repo.cfg.ElasticIndexName).
		Id(note.Id.String()).
		BodyJson(noteMap).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return ErrCantGetResponse
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) UpdateNote(ctx context.Context, note models.Note) error {
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

	start := time.Now()
	_, err = repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(note.Id.String()).
		Doc(noteMap).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return ErrCantGetResponse
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.elastic.Delete().
		Index(repo.cfg.ElasticIndexName).
		Id(id.String()).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return ErrCantGetResponse
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) AddSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.children.add(params.childID)").
		Lang("painless").
		Param("childID", childID.String())

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(id.String()).
		Script(script).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) RemoveSubNote(ctx context.Context, id uuid.UUID, childID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.children.removeIfContains(params.childID)").
		Lang("painless").
		Param("childID", childID.String())

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(id.String()).
		Script(script).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) AddCollaborator(ctx context.Context, noteID uuid.UUID, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.collaborators.add(params.collaboratorID)").
		Lang("painless").
		Param("collaboratorID", userID.String())

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
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

	script := elastic.NewScript("ctx._source.tags.add(params.tagName)").
		Lang("painless").
		Param("tagName", tagName)

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) DeleteTag(ctx context.Context, tagName string, noteID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("if (ctx._source.tags.contains(params.tagName)) { ctx._source.tags.remove(ctx._source.tags.indexOf(params.tagName)) }").
		Lang("painless").
		Param("tagName", tagName)

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Script(script).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) DeleteTagFromAllNotes(ctx context.Context, tagName string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.tags.remove(ctx._source.tags.indexOf(params.tagName))").
		Lang("painless").
		Param("tagName", tagName)

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("owner_id", userID)).
		Must(elastic.NewTermQuery("tags", tagName))

	searchResult, err := repo.elastic.Search().
		Index(repo.cfg.ElasticIndexName).
		Query(query).
		Do(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bulkRequest := repo.elastic.Bulk()
	for _, hit := range searchResult.Hits.Hits {
		bulkRequest = bulkRequest.Add(elastic.NewBulkUpdateRequest().
			Index(repo.cfg.ElasticIndexName).
			Id(hit.Id).
			Script(script),
		)
	}

	start := time.Now()
	_, err = bulkRequest.Do(context.Background())
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) UpdateTagOnAllNotes(ctx context.Context, oldTag string, newTag string, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	script := elastic.NewScript("ctx._source.tags[ctx._source.tags.indexOf(params.oldTag)] = params.NewTag").
		Lang("painless").
		Params(map[string]interface{}{
			"oldTag": oldTag,
			"newTag": newTag,
		})

	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("owner_id", userID)).
		Must(elastic.NewTermQuery("tags", oldTag))

	searchResult, err := repo.elastic.Search().
		Index(repo.cfg.ElasticIndexName).
		Query(query).
		Do(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bulkRequest := repo.elastic.Bulk()
	for _, hit := range searchResult.Hits.Hits {
		bulkRequest = bulkRequest.Add(elastic.NewBulkUpdateRequest().
			Index(repo.cfg.ElasticIndexName).
			Id(hit.Id).
			Script(script),
		)
	}

	start := time.Now()
	_, err = bulkRequest.Do(context.Background())
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) SetIcon(ctx context.Context, noteID uuid.UUID, icon string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Doc(map[string]interface{}{"icon": icon}).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) SetHeader(ctx context.Context, noteID uuid.UUID, header string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Doc(map[string]interface{}{"header": header}).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) ChangeFlag(ctx context.Context, noteID uuid.UUID, flag bool) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Doc(map[string]interface{}{"favorite": flag}).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) SetPublic(ctx context.Context, noteID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Doc(map[string]interface{}{"public": true}).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) SetPrivate(ctx context.Context, noteID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	start := time.Now()
	_, err := repo.elastic.Update().
		Index(repo.cfg.ElasticIndexName).
		Id(noteID.String()).
		Doc(map[string]interface{}{"public": false}).
		Do(ctx)
	repo.metr.ObserveResponseTime(log.GFN(), time.Since(start).Seconds())
	if err != nil {
		logger.Error(err.Error())
		repo.metr.IncreaseErrors(log.GFN())
		return err
	}

	logger.Info("success")
	return nil
}
