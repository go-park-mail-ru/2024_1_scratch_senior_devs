package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizeNote(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "<script>alert('XSS attack')</script>",
			output: "&lt;script&gt;alert(&#39;XSS attack&#39;)&lt;/script&gt;",
		},
		{
			input:  "Hello, World!",
			output: "Hello, World!",
		},
		{
			input:  "<h1>Hello</h1>",
			output: "&lt;h1&gt;Hello&lt;/h1&gt;",
		},
	}

	for _, tt := range tests {
		note := Note{Data: tt.input}
		note.Sanitize()

		assert.Equal(t, tt.output, note.Data)
	}
}
