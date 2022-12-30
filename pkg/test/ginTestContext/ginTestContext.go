package ginTestContext

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	ctx.Request = &http.Request{
		Method: "GET",
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx, responseWriter
}

func ContextWithClaim(claimName string, claimValue string) (*gin.Context, *httptest.ResponseRecorder) {
	ctx, response := TestContext()
	claims := jwt.MapClaims{}
	claims[claimName] = claimValue
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ctx.Set("token", token)
	return ctx, response
}

func ContextWithClaims(claims jwt.MapClaims) (*gin.Context, *httptest.ResponseRecorder) {
	ctx, response := TestContext()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ctx.Set("token", token)
	return ctx, response
}
