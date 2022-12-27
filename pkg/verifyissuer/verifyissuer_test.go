package verifyissuer

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"testing/quick"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func doTestRequest(claimName string, acceptableClaimValue string, actualClaimValue string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	claims := jwt.MapClaims{}
	claims[claimName] = actualClaimValue

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ctx.Request.Method = "GET"
	ctx.Set("token", token)

	claimsSelector := func(c *gin.Context) jwt.MapClaims {
		value, found := c.Get("token")
		if !found {
			panic("Token not found in context")
		}
		return value.(*jwt.Token).Claims.(jwt.MapClaims)
	}
	Verify(acceptableClaimValue, claimsSelector)(ctx)
	return ctx, response
}

func TestNonMatchingValues(t *testing.T) {
	test := func(authorizedClaim string, actualClaim string) bool {
		if authorizedClaim == actualClaim {
			actualClaim = actualClaim + "X"
		}
		_, response := doTestRequest("iss", authorizedClaim, actualClaim)
		return response.Code == 401
	}

	c := quick.Config{MaxCount: 100000}
	if err := quick.Check(test, &c); err != nil {
		t.Error(err)
	}
}
func TestMatchingValues(t *testing.T) {
	test := func(actualClaim string) bool {
		_, response := doTestRequest("iss", actualClaim, actualClaim)
		return response.Code == 200
	}

	c := quick.Config{MaxCount: 100000}
	if err := quick.Check(test, &c); err != nil {
		t.Error(err)
	}
}
