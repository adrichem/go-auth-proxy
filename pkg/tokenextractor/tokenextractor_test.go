package tokenextractor

import (
	"net/http"
	"testing"
)

const AuthHeader = "Authorization"

func TestMissingHeader(t *testing.T) {
	var x http.Request
	_, err := ExtractTokenFromHttpRequest(&x)
	if err == nil {
		t.Fatal("Missing auth header did not cause error")
	}
}

func TestHeaderHasNoValue(t *testing.T) {
	var x http.Request
	headers := make(map[string][]string)
	headers[AuthHeader] = []string{}
	x.Header = headers
	_, err := ExtractTokenFromHttpRequest(&x)
	if err == nil {
		t.Fatal("Empty auth header value did not cause error")
	}
}

func TestHeaderHasEmptyString(t *testing.T) {
	var x http.Request
	headers := make(map[string][]string)
	headers[AuthHeader] = []string{""}
	x.Header = headers
	_, err := ExtractTokenFromHttpRequest(&x)
	if err == nil {
		t.Fatal("Empty auth header value did not cause error")
	}
}

func TestHeaderHasInvalidFormat(t *testing.T) {
	var x http.Request
	headers := make(map[string][]string)
	headers[AuthHeader] = []string{"this is not in format Bearer ..."}
	x.Header = headers
	_, err := ExtractTokenFromHttpRequest(&x)
	if err == nil {
		t.Fatal("Invalid auth header format did not cause error")
	}
}

func TestHeaderHasValidFormat(t *testing.T) {
	var x http.Request
	headers := make(map[string][]string)
	headers[AuthHeader] = []string{"Bearer thetokenstring"}
	x.Header = headers
	tokenString, err := ExtractTokenFromHttpRequest(&x)
	if err != nil {
		t.Fatal("Valid auth header shoud not throw error")
	}
	if tokenString != "thetokenstring" {
		t.Fatal("Did not get the token string from the auth header")
	}
}
