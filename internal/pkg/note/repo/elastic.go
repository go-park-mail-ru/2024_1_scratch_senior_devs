package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/elastic"
	"log/slog"
	"strings"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
)

type NoteRepo struct {
	elastic *elasticsearch.Client
	logger  *slog.Logger
}

func CreateNoteRepo(elastic *elasticsearch.Client, logger *slog.Logger) *NoteRepo {
	return &NoteRepo{
		elastic: elastic,
		logger:  logger,
	}
}

func (repo *NoteRepo) ReadAllNotes(ctx context.Context, userID uuid.UUID, count int64, offset int64, searchValue string) ([]models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	query := elastic.ReadAllNotesQuery(userID, count, offset, searchValue)
	res, err := repo.elastic.Search(
		repo.elastic.Search.WithContext(ctx),
		repo.elastic.Search.WithIndex(config.ElasticIndexName),
		repo.elastic.Search.WithBody(bytes.NewReader(query)),
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("response error: " + res.Status())
		return nil, fmt.Errorf("response error: %s", res.String())
	}

	notes, err := elastic.GetSearchedNotesFromResponse(res)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("success")
	return notes, nil
}

func (repo *NoteRepo) ReadNote(ctx context.Context, noteId uuid.UUID) (models.Note, error) {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	res, err := repo.elastic.Get(config.ElasticIndexName, noteId.String())
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("response error: " + res.Status())
		return models.Note{}, fmt.Errorf("response error: %s", res.String())
	}

	resultNote, err := elastic.GetNoteFromResponse(res)
	if err != nil {
		logger.Error(err.Error())
		return models.Note{}, err
	}

	logger.Info("success")
	return resultNote, nil
}

func (repo *NoteRepo) CreateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	elasticNote, err := elastic.GetElasticNote(note)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	noteJSON, err := json.Marshal(elasticNote)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	req := esapi.IndexRequest{
		Index:      config.ElasticIndexName,
		DocumentID: elasticNote.Id.String(),
		Body:       bytes.NewReader(noteJSON),
	}

	res, err := req.Do(ctx, repo.elastic)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("response error: " + res.Status())
		return fmt.Errorf("response error: %s", res.Status())
	}

	logger.Info("success")
	return nil
}

func (repo *NoteRepo) UpdateNote(ctx context.Context, note models.Note) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	elasticNote, err := elastic.GetElasticNote(note)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	noteJSON, err := json.Marshal(models.ElasticNoteUpdate{Doc: elasticNote})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	fmt.Println(string(noteJSON))
	fmt.Println(elasticNote.Id.String())

	res, err := repo.elastic.Update(config.ElasticIndexName, elasticNote.Id.String(), strings.NewReader(string(noteJSON)))
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("response error: " + res.Status())
		return fmt.Errorf("response error: %s", res.Status())
	}

	logger.Info("success")
	return nil

}

func (repo *NoteRepo) DeleteNote(ctx context.Context, id uuid.UUID) error {
	logger := repo.logger.With(slog.String("ID", log.GetRequestId(ctx)), slog.String("func", log.GFN()))

	res, err := repo.elastic.Delete(config.ElasticIndexName, id.String())
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("response error: " + res.Status())
		return fmt.Errorf("response error: %s", res.Status())
	}

	logger.Info("success")
	return nil
}
