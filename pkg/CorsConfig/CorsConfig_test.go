package CorsConfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingfile(t *testing.T) {
	assert.Panics(t, func() { ReadConfig("") })
}

func TestInvalidFile(t *testing.T) {
	assert.Panics(t, func() { ReadConfig("invalid_cors.json") })
}
