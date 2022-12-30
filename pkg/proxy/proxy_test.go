package proxy

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
func contextWithGinResponseRecorder() (*gin.Context, *TestResponseRecorder) {
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

func TestBadUrlShouldPanic(t *testing.T) {
	assert.Panics(t, func() { Proxy("\n", "", "") })
	assert.Panics(t, func() { Proxy("http://", "", "") })
	assert.Panics(t, func() { Proxy("localhost", "", "") })
}

func TestHeaders(t *testing.T) {
	headerTest := func(headerName, expectedHeader string) bool {
		keyFound := false
		actualHeader := ""
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, keyFound = r.Header[headerName]
			actualHeader = r.Header.Get(headerName)
			w.WriteHeader(200)
		})

		srv := httptest.NewServer(mux)
		fnProxy := Proxy("http://"+srv.Listener.Addr().String(), headerName, expectedHeader)
		c, _ := contextWithGinResponseRecorder()
		fnProxy(c)
		srv.Close()
		return keyFound && actualHeader == expectedHeader
	}

	randomString := func(alphabet string, size int, r *rand.Rand) string {
		var buffer bytes.Buffer
		for i := 0; i < size; i++ {
			index := r.Intn(len(alphabet))
			buffer.WriteString(string(alphabet[index]))
		}
		return buffer.String()
	}

	RandHttpHeaderNameString := func(r *rand.Rand) string {
		alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
		return textproto.CanonicalMIMEHeaderKey(randomString(alphabet, 1+r.Intn(254), r))
	}

	RandHttpValueNameString := func(r *rand.Rand) string {
		alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_:;.,\\/\"'?!(){}[]@<>=-+*#$&`|~^%"
		return randomString(alphabet, r.Intn(255), r)
	}

	c := quick.Config{MaxCount: 1000, Values: func(values []reflect.Value, r *rand.Rand) {
		values[0] = reflect.ValueOf(RandHttpHeaderNameString(r))
		values[1] = reflect.ValueOf(RandHttpValueNameString(r))
	}}
	if err := quick.Check(headerTest, &c); err != nil {
		t.Error(err)
	}
}
