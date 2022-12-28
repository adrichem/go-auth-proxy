package claimsSelector

import (
	"go-auth-proxy/pkg/test/ginTestContext"
	"testing"
)

func TestTokenMissing(t *testing.T) {
	defer func() {
		recover()
	}()
	ctx, _ := ginTestContext.TestContext()
	FromGinContext("token")(ctx)
	t.Fatalf("Expected a panic, non found")
}

func TestTokenInvalidType(t *testing.T) {
	defer func() {
		recover()
	}()
	ctx, _ := ginTestContext.TestContext()
	ctx.Set("token", "I AM NOT A MAPCLAIMS STRUCT")
	FromGinContext("token")(ctx)
	t.Fatalf("Expected a panic, non found")
}
