package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGFN(t *testing.T) {
	assert.Equal(t, "utils.TestGFN", GFN())
}
