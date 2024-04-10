package repo

import (
	"log/slog"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

var testElasticLogger *slog.Logger

func init() {
	testElasticLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func TestNoteRepo_MakeHelloNoteData(t *testing.T) {
	testConfig := config.ElasticConfig{
		ElasticIndexName:            "note",
		ElasticSearchValueMinLength: 2,
	}
	username := "testuser"
	expected := []byte(`
		{
			"title": "You-note ❤️",
			"content": [
				{
					"id": "1",
					"type": "div",
					"content": [
						{
							"id": "2",
							"content": "Привет, testuser!"
						}
					]
				}
			]
		}
	`)

	t.Run("Test_MakeHelloNoteData", func(t *testing.T) {
		repo := CreateNoteElastic(&elastic.Client{}, testElasticLogger, testConfig)
		result := repo.MakeHelloNoteData(username)

		assert.Equal(t, expected, result)
	})
}
