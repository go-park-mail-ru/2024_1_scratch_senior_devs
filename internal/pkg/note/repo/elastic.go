package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
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
	cfg     config.ElasticConfig
}

func CreateNoteRepo(elastic *elastic.Client, logger *slog.Logger, cfg config.ElasticConfig) *NoteRepo {
	return &NoteRepo{
		elastic: elastic,
		logger:  logger,
		cfg:     cfg,
	}
}

func (repo *NoteRepo) ReadAllNotes(ctx context.Context, userID uuid.UUID, count int64, offset int64, searchValue string) ([]models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	userIdQuery := elastic.NewTermsQuery("owner_id", strings.ToLower(userID.String()))
	searchQuery := elastic.NewMatchQuery("data", searchValue)

	totalQuery := repo.elastic.Search().Query(elastic.NewBoolQuery().Must(userIdQuery, searchQuery))
	if utf8.RuneCountInString(searchValue) < repo.cfg.ElasticSearchValueMinLength {
		totalQuery = repo.elastic.Search().Query(userIdQuery)
	}

	search, err := totalQuery.
		Index(repo.cfg.ElasticIndexName).
		From(int(offset)).
		Size(int(count)).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return []models.Note{}, errors.New("can`t get response")
	}

	notes := make([]models.Note, 0)
	for _, hit := range search.Hits.Hits {
		note := models.ElasticNote{}
		if err := json.Unmarshal(hit.Source, &note); err != nil {
			logger.Error(err.Error())
			return []models.Note{}, err
		}
		notes = append(notes, models.Note{
			Id:         note.Id,
			Data:       []byte(note.Data),
			CreateTime: note.CreateTime,
			UpdateTime: note.UpdateTime,
			OwnerId:    note.OwnerId,
		})
	}

	logger.Info("success")
	return notes, nil
}

func (repo *NoteRepo) ReadNote(ctx context.Context, noteID uuid.UUID) (models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	search, err := repo.elastic.Search().
		Index(repo.cfg.ElasticIndexName).
		Query(elastic.NewTermQuery("_id", noteID)).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, errors.New("can`t get response")
	}

	if len(search.Hits.Hits) == 0 {
		logger.Error("note not found")
		return models.Note{}, errors.New("note not found")
	}

	note := models.ElasticNote{}
	if err := json.Unmarshal(search.Hits.Hits[0].Source, &note); err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return models.Note{
		Id:         note.Id,
		Data:       []byte(note.Data),
		CreateTime: note.CreateTime,
		UpdateTime: note.UpdateTime,
		OwnerId:    note.OwnerId,
	}, nil
}

func (repo *NoteRepo) CreateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	elasticNote := models.ElasticNote{
		Id:         note.Id,
		Data:       string(note.Data),
		CreateTime: note.CreateTime,
		UpdateTime: note.UpdateTime,
		OwnerId:    note.OwnerId,
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
		Index(repo.cfg.ElasticIndexName).
		Id(elasticNote.Id.String()).
		BodyJson(noteMap).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("can`t get response")
	}

	logger.Info("success")
	return nil
}

func (repo *NoteRepo) UpdateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	elasticNote := models.ElasticNote{
		Id:         note.Id,
		Data:       string(note.Data),
		CreateTime: note.CreateTime,
		UpdateTime: note.UpdateTime,
		OwnerId:    note.OwnerId,
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
		Index(repo.cfg.ElasticIndexName).
		Id(elasticNote.Id.String()).
		Doc(noteMap).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("can`t get response")
	}

	logger.Info("success")
	return nil

}

func (repo *NoteRepo) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	_, err := repo.elastic.Delete().
		Index(repo.cfg.ElasticIndexName).
		Id(id.String()).
		Do(ctx)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("can`t get response")
	}

	logger.Info("success")
	return nil
}

func (repo *NoteRepo) MakeHelloNoteData(username string) []byte {
	return []byte(fmt.Sprintf(`
		{
			"title": "You-note ❤️",
			"content": [
				{
					"id": "1",
					"type": "div",
					"content": [
						{
							"id": "2",
							"content": "Привет, %s!"
						}
					]
				}
			]
		}
	`, username))
}
