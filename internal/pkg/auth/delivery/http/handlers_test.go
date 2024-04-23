package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHelloNoteData(t *testing.T) {
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
		result := makeHelloNoteData(username)

		assert.Equal(t, expected, result)
	})
}
