package ginTestContext

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func (r *TestResponseRecorder) closeClient() {
	r.closeChannel <- true
}

func CreateTestResponseRecorder() *TestResponseRecorder {
	return &TestResponseRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func TestContext() (*gin.Context, *TestResponseRecorder) {
	gin.SetMode(gin.TestMode)
	responseWriter := CreateTestResponseRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)
	ctx.Request = &http.Request{
		Method: "GET",
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx, responseWriter
}

func ContextWithClaim(claimName string, claimValue string) (*gin.Context, *TestResponseRecorder) {
	ctx, response := TestContext()
	claims := jwt.MapClaims{}
	claims[claimName] = claimValue
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ctx.Set("token", token)
	return ctx, response
}

func ContextWithClaims(claims jwt.MapClaims) (*gin.Context, *TestResponseRecorder) {
	ctx, response := TestContext()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ctx.Set("token", token)
	return ctx, response
}
