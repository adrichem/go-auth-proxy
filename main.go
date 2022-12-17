package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	var ListenAddress = flag.String("address", ":80", "Adress to listen on.")
	var Upstream = flag.String("upstream", "", "Url to proxy to.")
	var HeaderName = flag.String("header-name", "apikey", "Name of header to add to proxied requests.")
	var HeaderValue = flag.String("header-value", "", "Value of header to add to proxied requests.")
	flag.Parse()
	if *ListenAddress == "" || *Upstream == "" || *HeaderName == "" || *HeaderValue == "" {
		flag.PrintDefaults()
		panic("invalid arguments")
	}
	http.HandleFunc("/", authenticateAzureAd(proxy(*Upstream, *HeaderName, *HeaderValue)))
	fmt.Printf("Listening on %s\n", *ListenAddress)
	log.Fatal(http.ListenAndServe(*ListenAddress, nil))
}

func proxy(Upstream string, HeaderName string, HeaderValue string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		remote, err := url.Parse(Upstream)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			r.Header.Add(HeaderName, HeaderValue)
			req.Header = r.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = r.RequestURI
		}
		proxy.ServeHTTP(w, r)

	}
}

func authenticateAzureAd(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := verifyToken(r)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		next(w, r)
	}
}

func extractToken(r *http.Request) (string, error) {
	authHeaderArray, ok := r.Header["Authorization"]
	if !ok {
		return "", errors.New("authorization header is missing")
	}

	if len(authHeaderArray) == 0 {
		return "", errors.New("authorization header is empty")
	}
	authHeader := authHeaderArray[0]
	words := strings.Split(authHeader, " ")
	if len(words) < 2 || words[0] != "Bearer" {
		return "", errors.New("authorization is not in format 'Bearer ...'")
	}
	return words[1], nil
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString, err := extractToken(r)
	if err != nil {
		return nil, err
	}
	//Rerieve the list of public keys from Azure AD
	jwks, err := keyfunc.Get("https://login.microsoftonline.com/common/discovery/v2.0/keys", keyfunc.Options{})
	if err != nil {
		return nil, err
	}
	//Validate and parse the JWT token
	return jwt.Parse(tokenString, jwks.Keyfunc)
}
