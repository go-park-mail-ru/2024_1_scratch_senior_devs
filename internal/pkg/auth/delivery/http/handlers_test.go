package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHelloNoteData(t *testing.T) {
	username := "testuser"
	expected := []byte(`
	{
		"title": "",
		"content": [
		    {
			   "pluginName": "textBlock",
			   "content": "Hello testuser"
		    },
		    {
			   "pluginName": "div",
			   "children": [
				  {
					 "pluginName": "br",
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
