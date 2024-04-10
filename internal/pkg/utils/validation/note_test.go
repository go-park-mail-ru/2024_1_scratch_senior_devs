package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckNoteTitle(t *testing.T) {
	tests := []struct {
		name     string
		noteData []byte
		isErr    bool
	}{
		{
			name:     "Test_CheckNoteTitle_Success",
			noteData: []byte(`{"title": "my note"}`),
			isErr:    false,
		},
		{
			name:     "Test_CheckNoteTitle_Fail_1",
			noteData: []byte(`{"title": "my note"`),
			isErr:    true,
		},
		{
			name:     "Test_CheckNoteTitle_Fail_2",
			noteData: []byte(`{"titl": "my note"}`),
			isErr:    true,
		},
		{
			name:     "Test_CheckNoteTitle_Fail_3",
			noteData: []byte(`{"title": ["abc", "def"]}`),
			isErr:    true,
		},
		{
			name:     "Test_CheckNoteTitle_Fail_4",
			noteData: []byte(`{"title": ""}`),
			isErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckNoteTitle(tt.noteData)

			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
