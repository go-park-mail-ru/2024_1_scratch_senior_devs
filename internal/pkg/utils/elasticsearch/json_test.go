package elasticsearch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessValue(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		result interface{}
	}{
		{
			name: "Test_ProcessValue_Success_1",
			value: map[string]interface{}{
				"title":   "my note",
				"content": "some text",
			},
			result: []map[string]interface{}{
				{
					"key":   "title",
					"value": "my note",
				},
				{
					"key":   "content",
					"value": "some text",
				},
			},
		},
		{
			name:   "Test_ProcessValue_Success_2",
			value:  []interface{}{"abc", "def"},
			result: []interface{}{"abc", "def"},
		},
		{
			name:   "Test_ProcessValue_Success_3",
			value:  "abc",
			result: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processValue(tt.value)

			assert.Equal(t, result, tt.result)
		})
	}
}
