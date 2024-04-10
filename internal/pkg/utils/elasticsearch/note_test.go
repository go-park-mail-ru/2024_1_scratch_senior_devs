package elasticsearch

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertToElasticNote(t *testing.T) {
	id := uuid.NewV4()
	currTime := time.Now().UTC()

	tests := []struct {
		name   string
		note   models.Note
		result models.ElasticNote
		isErr  bool
	}{
		{
			name: "Test_ConvertToElasticNote_Success",
			note: models.Note{
				Id:         id,
				Data:       []byte(`{"title": "my note"}`),
				CreateTime: currTime,
				UpdateTime: currTime,
				OwnerId:    id,
			},
			result: models.ElasticNote{
				Id:         id,
				Data:       `{"title": "my note"}`,
				CreateTime: currTime,
				UpdateTime: currTime,
				OwnerId:    id,
			},
			isErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elasticNote := ConvertToElasticNote(tt.note)

			assert.Equal(t, elasticNote, tt.result)
		})
	}
}

func TestConvertToUsualNote(t *testing.T) {
	id := uuid.NewV4()
	currTime := time.Now().UTC()

	tests := []struct {
		name   string
		note   models.ElasticNote
		result models.Note
	}{
		{
			name: "Test_ConvertToUsualNote_Success",
			note: models.ElasticNote{
				Id:         id,
				Data:       `{"title": "my note"}`,
				CreateTime: currTime,
				UpdateTime: currTime,
				OwnerId:    id,
			},
			result: models.Note{
				Id:         id,
				Data:       []byte(`{"title": "my note"}`),
				CreateTime: currTime,
				UpdateTime: currTime,
				OwnerId:    id,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := ConvertToUsualNote(tt.note)

			assert.Equal(t, note, tt.result)
		})
	}
}
