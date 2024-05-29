package exportpdf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareHTML(t *testing.T) {
	_, _, _, err := prepareHTML(TestNoteHTMLInput)
	assert.Nil(t, err)
}
