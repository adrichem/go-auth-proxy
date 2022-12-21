package main

import (
	"flag"
	"go-auth-proxy/pkg/tokenextractor"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	var ListenAddress = flag.String("address", ":80", "Adress to listen on.")
	var Upstream = flag.String("upstream", "", "Url to proxy to.")
	var HeaderName = flag.String("header-name", "Apikey", "Name of header to add to proxied requests.")
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
	r.Use(AzureAdJwtTokenValidation())
	if *Upstream != "" {
		//Proxying to upstream
		r.Any("/*path", proxy(*Upstream, *HeaderName, *HeaderValue))
	} else {
		//Load test mode. Just return HTTP 200
		r.Any("/*path", func(c *gin.Context) { c.Status(http.StatusOK) })
	}
	err := r.Run(*ListenAddress)
	if err != nil {
		panic(err)
	}
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

func AzureAdJwtTokenValidation() gin.HandlerFunc {
	//Azure AD keyset is independent from the token content
	jwks, err := keyfunc.Get("https://login.microsoftonline.com/common/discovery/v2.0/keys", keyfunc.Options{})
	if err != nil {
		panic(err)
	}
	keyFuncSelector := func(string) (*keyfunc.JWKS, error) { return jwks, nil }
	return createAuthenticationMiddleware(keyFuncSelector, tokenextractor.ExtractTokenFromHttpRequest)
}

func createAuthenticationMiddleware(selectKeySet func(string) (*keyfunc.JWKS, error),
	extractToken func(r *http.Request) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var err error
		var jwks *keyfunc.JWKS

		token, err = extractToken(c.Request)
		if err == nil {
			jwks, err = selectKeySet(token)
		}
		if err == nil {
			_, err = jwt.Parse(token, jwks.Keyfunc)
		}
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
