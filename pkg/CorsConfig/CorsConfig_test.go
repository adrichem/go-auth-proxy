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

func TestValidFile(t *testing.T) {
	cors := ReadConfig("valid_cors.json")
	assert.Len(t, cors.AllowOrigins, 1)
	assert.Contains(t, cors.AllowOrigins, "*")

	assert.Len(t, cors.AllowHeaders, 2)
	assert.Contains(t, cors.AllowHeaders, "Authorization")
	assert.Contains(t, cors.AllowHeaders, "Origin")

	assert.Len(t, cors.AllowMethods, 3)
	assert.Contains(t, cors.AllowMethods, "GET")
	assert.Contains(t, cors.AllowMethods, "PUT")
	assert.Contains(t, cors.AllowMethods, "POST")

}
