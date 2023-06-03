package claimsSelector

import (
	"testing"

	"github.com/adrichem/go-auth-proxy/pkg/test/ginTestContext"

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
