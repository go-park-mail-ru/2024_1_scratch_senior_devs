package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/satori/uuid"
	"unicode/utf8"
)

func ReadAllNotesQuery(userID uuid.UUID, count int64, offset int64, searchValue string) []byte {
	if utf8.RuneCountInString(searchValue) < config.ElasticSearchValueMinLength {
		return []byte(fmt.Sprintf(`
			{
			  "query": {
				"match": {
				  "owner_id": "%s"
				}
			  },
			  "from": %d,
			  "size": %d
			}
		`, userID, offset, count))
	}

	return []byte(fmt.Sprintf(`
		{
		  "query": {
			"bool": {
			  "must": [
				{
				  "nested": {
					"path": "elastic_data",
					"query": {
					  "match": {
						"elastic_data.value": {
						  "query": "%s",
						  "operator": "and"
						}
					  }
					}
				  }
				},
				{
				  "term": {
					"owner_id": {
					  "value": "%s"
					}
				  }
				}
			  ]
			}
		  },
		  "from": %d,
		  "size": %d
		}
	`, searchValue, userID.String(), offset, count))
}

func GetSearchedNotesFromResponse(res *esapi.Response) ([]models.Note, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	notes := make([]models.Note, 0)

	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, errors.New("incorrect response")
	}

	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			return nil, errors.New("incorrect response")
		}

		elasticNoteInterface, err := json.Marshal(source)
		if err != nil {
			return nil, errors.New("incorrect response")
		}

		elasticNote := models.ElasticNote{}
		if err := json.Unmarshal(elasticNoteInterface, &elasticNote); err != nil {
			return nil, errors.New("incorrect response")
		}

		notes = append(notes, models.Note{
			Id:         elasticNote.Id,
			Data:       []byte(elasticNote.Data),
			CreateTime: elasticNote.CreateTime,
			UpdateTime: elasticNote.UpdateTime,
			OwnerId:    elasticNote.OwnerId,
		})
	}

	return notes, nil
}

func GetNoteFromResponse(res *esapi.Response) (models.Note, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return models.Note{}, err
	}

	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return models.Note{}, errors.New("incorrect response")
	}

	elasticNoteInterface, err := json.Marshal(source)
	if err != nil {
		return models.Note{}, errors.New("incorrect response")
	}

	elasticNote := models.ElasticNote{}
	if err := json.Unmarshal(elasticNoteInterface, &elasticNote); err != nil {
		return models.Note{}, errors.New("incorrect response")
	}

	return models.Note{
		Id:         elasticNote.Id,
		Data:       []byte(elasticNote.Data),
		CreateTime: elasticNote.CreateTime,
		UpdateTime: elasticNote.UpdateTime,
		OwnerId:    elasticNote.OwnerId,
	}, nil
}
