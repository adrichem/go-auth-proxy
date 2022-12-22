package main

import (
	"errors"
	"flag"
	"go-auth-proxy/pkg/tokenextractor"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func getParsedToken(c *gin.Context, paramName string) *jwt.Token {
	value, found := c.Get(paramName)
	if !found {
		panic("Token not found in context")
	}
	return value.(*jwt.Token)
}

func fail(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, err.Error())
	c.Abort()
}

func main() {
	var ListenAddress = flag.String("address", ":80", "Adress to listen on.")
	var Upstream = flag.String("upstream", "", "Url to proxy to.")
	var HeaderName = flag.String("header-name", "Apikey", "Name of header to add to proxied requests.")
	var HeaderValue = flag.String("header-value", "", "Value of header to add to proxied requests.")
	var Aud = flag.String("aud", "", "Reject tokens not issued to intendend audience (aud claim)")
	var Iss = flag.String("iss", "", "Reject tokens not issued by intended issuer (iss claim")
	flag.Parse()
	if *ListenAddress == "" {
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
	r.Use(azureAdJwtTokenValidation("token"))
	r.Use(verifyAudience(*Aud, "token"))
	r.Use(verifyIssuer(*Iss, "token"))
	if *Upstream != "" {
		log.Printf("Proxying to %s", *Upstream)
		r.Any("/*path", proxy(*Upstream, *HeaderName, *HeaderValue))
	} else {
		log.Println("Running in load test mode. All requests get HTTP 200 response")
		r.Any("/*path", func(c *gin.Context) { c.Status(http.StatusOK) })
	}

	log.Printf("Listening on %s", *ListenAddress)
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

func azureAdJwtTokenValidation(paramName string) gin.HandlerFunc {
	//Azure AD keyset is independent from the token content
	jwks, err := keyfunc.Get("https://login.microsoftonline.com/common/discovery/v2.0/keys", keyfunc.Options{})
	if err != nil {
		panic(err)
	}
	keyFuncSelector := func(string) (*keyfunc.JWKS, error) { return jwks, nil }
	return createAuthenticationMiddleware(keyFuncSelector, tokenextractor.ExtractTokenFromHttpRequest, paramName)
}

func createAuthenticationMiddleware(selectKeySet func(string) (*keyfunc.JWKS, error),
	extractToken func(r *http.Request) (string, error),
	paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var err error
		var jwks *keyfunc.JWKS
		var parsedToken *jwt.Token
		token, err = extractToken(c.Request)
		if err == nil {
			jwks, err = selectKeySet(token)
		}
		if err == nil {
			parsedToken, err = jwt.Parse(token, jwks.Keyfunc)
		}
		if err == nil {
			c.Set(paramName, parsedToken)
		}
		if err != nil {
			fail(c, err)
			return
		}
		c.Next()
	}
}

func verifyClaim(tokenParamName string, failureMsg string, predicate func(jwt.MapClaims) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := getParsedToken(c, tokenParamName).Claims.(jwt.MapClaims)
		if predicate(claims) {
			c.Next()
		} else {
			fail(c, errors.New(failureMsg))
		}
	}
}

func verifyIssuer(ExpectedIssuer string, tokenParamName string) gin.HandlerFunc {
	return verifyClaim(tokenParamName, "invalid issuer", func(c jwt.MapClaims) bool { return ExpectedIssuer != "" && c.VerifyIssuer(ExpectedIssuer, true) })
}

func verifyAudience(ExpectedAudience string, tokenParamName string) gin.HandlerFunc {
	return verifyClaim(tokenParamName, "invalid audience", func(c jwt.MapClaims) bool { return ExpectedAudience != "" && c.VerifyIssuer(ExpectedAudience, true) })
}
