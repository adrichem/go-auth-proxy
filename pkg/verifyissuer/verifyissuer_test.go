package verifyissuer

import (
	"go-auth-proxy/pkg/claimsSelector"
	"go-auth-proxy/pkg/test/ginTestContext"
	"testing"
	"testing/quick"
)

func TestNonMatchingValues(t *testing.T) {
	test := func(authorizedClaim string, actualClaim string) bool {
		if authorizedClaim == actualClaim {
			actualClaim = actualClaim + "X"
		}
		ctx, response := ginTestContext.ContextWithClaim("iss", actualClaim)
		Verify(authorizedClaim, claimsSelector.FromGinContext("token"))(ctx)
		return response.Code == 401
	}

	c := quick.Config{MaxCount: 100000}
	if err := quick.Check(test, &c); err != nil {
		t.Error(err)
	}
}
func TestMatchingValues(t *testing.T) {
	test := func(actualClaim string) bool {
		ctx, response := ginTestContext.ContextWithClaim("iss", actualClaim)
		Verify(actualClaim, claimsSelector.FromGinContext("token"))(ctx)
		return response.Code == 200
	}

	c := quick.Config{MaxCount: 100000}
	if err := quick.Check(test, &c); err != nil {
		t.Error(err)
	}
}
