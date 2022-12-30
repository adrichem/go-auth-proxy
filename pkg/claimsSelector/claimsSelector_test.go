package claimsSelector

import (
	"go-auth-proxy/pkg/test/ginTestContext"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenMissing(t *testing.T) {
	assert.Panics(t, func() {
		ctx, _ := ginTestContext.TestContext()
		FromGinContext("token")(ctx)
	})
}

func TestTokenInvalidType(t *testing.T) {
	assert.Panics(t, func() {
		ctx, _ := ginTestContext.TestContext()
		ctx.Set("token", "I AM NOT A MAPCLAIMS STRUCT")
		FromGinContext("token")(ctx)
	})
}
