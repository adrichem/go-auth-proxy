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

/*
When testing the combination of gin and the http reverse proxy, we get panics like:
panic: interface conversion: *httptest.ResponseRecorder is not http.CloseNotifier: missing method CloseNotify
For this we need a have a custom recorder that implements httptest.ResponseRecorder and the CloseNotify method
*/
func TestContextWithGinResponseRecorder() (*gin.Context, *TestResponseRecorder) {
	rr := CreateTestResponseRecorder()
	// creates a test context and gin engine
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = &http.Request{
		Method: "GET",
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx, rr
}

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
