package elasticsearch

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

func processValue(value interface{}) interface{} {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		outputSlice := make([]map[string]interface{}, 0)
		for key, val := range typedValue {
			outputSlice = append(outputSlice, map[string]interface{}{"key": key, "value": processValue(val)})
		}
		return outputSlice
	case []interface{}:
		outputSlice := make([]interface{}, 0)
		for _, val := range typedValue {
			outputSlice = append(outputSlice, processValue(val))
		}
		return outputSlice
	default:
		return value
	}
}

func getElasticData(inputJSON []byte) (interface{}, error) {
	var inputMap map[string]interface{}
	if err := json.Unmarshal(inputJSON, &inputMap); err != nil {
		return []byte{}, err
	}

	return processValue(inputMap), nil
}

func ConvertToElasticNote(note models.Note) (models.ElasticNote, error) {
	elasticData, err := getElasticData(note.Data)
	if err != nil {
		return models.ElasticNote{}, err
	}

	return models.ElasticNote{
		Id:          note.Id,
		Data:        string(note.Data),
		ElasticData: elasticData,
		CreateTime:  note.CreateTime,
		UpdateTime:  note.UpdateTime,
		OwnerId:     note.OwnerId,
	}, nil
}

func ConvertToUsualNote(note models.ElasticNote) models.Note {
	return models.Note{
		Id:         note.Id,
		Data:       []byte(note.Data),
		CreateTime: note.CreateTime,
		UpdateTime: note.UpdateTime,
		OwnerId:    note.OwnerId,
	}
}
