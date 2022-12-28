package proxy

import (
	"bytes"
	"go-auth-proxy/pkg/test/ginTestContext"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {

	headerTest := func(headerName, expectedHeader string) bool {
		keyFound := false
		actualHeader := ""
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, keyFound = r.Header[headerName]
			actualHeader = r.Header.Get(headerName)
			respondWithStatus(http.StatusOK)(w, r)
		})

		srv := httptest.NewServer(mux)
		fnProxy := Proxy("http://"+srv.Listener.Addr().String(), headerName, expectedHeader)
		c, _ := ginTestContext.TestContext()
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

func TestUpstreamUnreachable(t *testing.T) {
	fnProxy := Proxy("http://aaaaaaaaaaaaa127.0.0.1:1", "", "")
	c, response := ginTestContext.TestContext()
	fnProxy(c)
	result := response.Result()
	assert.Equal(t, http.StatusBadGateway, result.StatusCode, result.Body)
}

func TestGet(t *testing.T) {
	expectedStatus := http.StatusConflict
	mux := http.NewServeMux()
	mux.HandleFunc("/", respondWithStatus(expectedStatus))
	srv := httptest.NewServer(mux)
	addr := srv.Listener.Addr()

	fnProxy := Proxy("http://"+addr.String(), "", "")
	c, response := ginTestContext.TestContext()
	fnProxy(c)
	assert.Equal(t, expectedStatus, response.Code)

	srv.Close()
}

func respondWithStatus(status int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}
}
