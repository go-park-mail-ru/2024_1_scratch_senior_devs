package repo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/elasticsearch"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/olivere/elastic/v7"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
)

type NoteRepo struct {
	elastic *elastic.Client
	logger  *slog.Logger
}

func CreateNoteRepo(elastic *elastic.Client, logger *slog.Logger) *NoteRepo {
	return &NoteRepo{
		elastic: elastic,
		logger:  logger,
	}
}

func (repo *NoteRepo) ReadAllNotes(ctx context.Context, userID uuid.UUID, count int64, offset int64, searchValue string) ([]models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	userIdQuery := elastic.NewTermsQuery("owner_id", strings.ToLower(userID.String()))
	searchQuery := elastic.NewNestedQuery("elastic_data", elastic.NewMatchQuery("elastic_data.value", searchValue))

	totalQuery := repo.elastic.Search().Query(elastic.NewBoolQuery().Must(userIdQuery, searchQuery))
	if utf8.RuneCountInString(searchValue) < config.ElasticSearchValueMinLength {
		totalQuery = repo.elastic.Search().Query(userIdQuery)
	}

	search, err := totalQuery.
		Index(config.ElasticIndexName).
		From(int(offset)).
		Size(int(count)).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return []models.Note{}, err
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

func (repo *NoteRepo) ReadNote(ctx context.Context, noteID uuid.UUID) (models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	search, err := repo.elastic.Search().
		Index(config.ElasticIndexName).
		Query(elastic.NewTermQuery("_id", noteID)).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	for _, hit := range search.Hits.Hits {
		note := models.ElasticNote{}
		if err := json.Unmarshal(hit.Source, &note); err != nil {
			logger.Error(err.Error())
			return models.Note{}, err
		}

		logger.Info("success")
		return elasticsearch.ConvertToUsualNote(note), nil
	}

	logger.Error("note not found")
	return models.Note{}, errors.New("note not found")
}

func (repo *NoteRepo) CreateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	elasticNote, err := elasticsearch.ConvertToElasticNote(note)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

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
		Index(config.ElasticIndexName).
		Id(elasticNote.Id.String()).
		BodyJson(noteMap).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}

func (repo *NoteRepo) UpdateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	elasticNote, err := elasticsearch.ConvertToElasticNote(note)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

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
		Index(config.ElasticIndexName).
		Id(elasticNote.Id.String()).
		Doc(noteMap).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil

}

func (repo *NoteRepo) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	_, err := repo.elastic.Delete().
		Index(config.ElasticIndexName).
		Id(id.String()).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("success")
	return nil
}
