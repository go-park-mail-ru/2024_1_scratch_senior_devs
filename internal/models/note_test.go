package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizeNote(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
	}{
		{
			input:  []byte("<script>alert('XSS attack')</script>"),
			output: []byte("&lt;script&gt;alert(&#39;XSS attack&#39;)&lt;/script&gt;"),
		},
		{
			input:  []byte("Hello, World!"),
			output: []byte("Hello, World!"),
		},
		{
			input:  []byte("<h1>Hello</h1>"),
			output: []byte("&lt;h1&gt;Hello&lt;/h1&gt;"),
		},
	}

	for _, tt := range tests {
		result := Sanitize(tt.input)

		assert.Equal(t, tt.output, result)
	}
}
