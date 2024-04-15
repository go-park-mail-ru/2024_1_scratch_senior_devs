package repo

import (
	"context"
	"encoding/json"
	"errors"
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

func (repo *NoteElastic) SearchNotes(ctx context.Context, userID uuid.UUID, count int64, offset int64, searchValue string) ([]models.Note, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	userIdQuery := elastic.NewTermsQuery("owner_id", strings.ToLower(userID.String()))
	searchQuery := elastic.NewMatchQuery("data", searchValue)

	search, err := repo.elastic.Search().
		Query(elastic.NewBoolQuery().Must(userIdQuery, searchQuery)).
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

func (repo *NoteElastic) CreateNote(ctx context.Context, note models.Note) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	elasticNote := elasticsearch.ConvertToElasticNote(note)

	noteJSON, err := json.Marshal(elasticNote)
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
		Id(elasticNote.Id.String()).
		BodyJson(noteMap).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return ErrCantGetResponse
	}

	//_, err = repo.elastic.Reindex().Do(ctx)
	//if err != nil {
	//	logger.Error(err.Error())
	//}

	logger.Info("success")
	return nil
}

func (repo *NoteElastic) UpdateNote(ctx context.Context, note models.Note) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GFN()))

	elasticNote := elasticsearch.ConvertToElasticNote(note)

	noteJSON, err := json.Marshal(elasticNote)
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
		Id(elasticNote.Id.String()).
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
