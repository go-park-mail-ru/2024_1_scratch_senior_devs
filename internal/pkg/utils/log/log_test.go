package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGFN(t *testing.T) {
	assert.Equal(t, "log.TestGFN", GFN())
}
