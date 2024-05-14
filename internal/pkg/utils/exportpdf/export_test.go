package exportpdf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePDF(t *testing.T) {
	_, _, err := GeneratePDF(TestNoteHTMLInput)
	assert.Nil(t, err)
}
