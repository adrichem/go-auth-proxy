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
			respondWithStatus(200)(w, r)
		})

		srv := httptest.NewServer(mux)
		fnProxy := Proxy("http://"+srv.Listener.Addr().String(), headerName, expectedHeader)
		c, _ := ginTestContext.TestContextWithGinResponseRecorder()
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

func respondWithStatus(status int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}
}
