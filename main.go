package main

import (
	"errors"
	"flag"
	"go-auth-proxy/pkg/claimverifier"
	jwtValidator "go-auth-proxy/pkg/jwtAuthenticationMiddleware"
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

func fail(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, err.Error())
	c.Abort()
}

func pass(c *gin.Context) { c.Next() }

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
	r.Use(verifyAudience(*Aud))
	r.Use(verifyIssuer(*Iss))
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
	//Azure AD keysetglobal and independant from token and AD tenant
	jwks, err := keyfunc.Get("https://login.microsoftonline.com/common/discovery/v2.0/keys", keyfunc.Options{})
	if err != nil {
		panic(err)
	}
	extractToken := func(c *gin.Context) (string, error) { return tokenextractor.ExtractTokenFromHttpRequest(c.Request) }
	pass := func(c *gin.Context, t *jwt.Token) { c.Set("token", t); pass(c) }
	return jwtValidator.Create(jwks.Keyfunc, extractToken, pass, fail)
}

func tokenFromGinContext(c *gin.Context) *jwt.Token {
	value, found := c.Get("token")
	if !found {
		panic("Token not found in context")
	}
	return value.(*jwt.Token)
}

func verifyIssuer(ExpectedIssuer string) gin.HandlerFunc {
	fnPredicate := func(c jwt.Claims) bool {
		return ExpectedIssuer != "" && c != nil && c.(jwt.MapClaims).VerifyIssuer(ExpectedIssuer, true)
	}
	fnFail := func(c *gin.Context) { fail(c, errors.New("invalid issuer")) }
	return claimverifier.VerifyClaim(tokenFromGinContext, fnPredicate, pass, fnFail)
}

func verifyAudience(ExpectedAudience string) gin.HandlerFunc {
	fnPredicate := func(c jwt.Claims) bool {
		return ExpectedAudience != "" && c != nil && c.(jwt.MapClaims).VerifyAudience(ExpectedAudience, true)
	}
	fnFail := func(c *gin.Context) { fail(c, errors.New("invalid audience")) }
	return claimverifier.VerifyClaim(tokenFromGinContext, fnPredicate, pass, fnFail)
}
