package main

import (
	"errors"
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	var ListenAddress = flag.String("address", "0.0.0.0:80", "Adress to listen on.")
	var Upstream = flag.String("upstream", "", "Url to proxy to.")
	var HeaderName = flag.String("header-name", "apikey", "Name of header to add to proxied requests.")
	var HeaderValue = flag.String("header-value", "", "Value of header to add to proxied requests.")
	flag.Parse()
	if *ListenAddress == "" || *Upstream == "" {
		flag.PrintDefaults()
		panic("invalid arguments")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Authorization"},
	}))
	r.Use(azureADAuthenticationMiddleware)
	r.Any("/*path", proxy(*Upstream, *HeaderName, *HeaderValue))
	r.Run(*ListenAddress)
}

func proxy(Upstream string, HeaderName string, HeaderValue string) gin.HandlerFunc {
	remote, err := url.Parse(Upstream)
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			if HeaderName != "" && HeaderValue != "" {
				req.Header.Add(HeaderName, HeaderValue)
			}
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("path")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func azureADAuthenticationMiddleware(c *gin.Context) {
	_, err := verifyToken(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}
	c.Next()
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
