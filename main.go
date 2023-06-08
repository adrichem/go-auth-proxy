package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/adrichem/go-auth-proxy/pkg/CorsConfig"
	"github.com/adrichem/go-auth-proxy/pkg/claimsSelector"
	jwtValidator "github.com/adrichem/go-auth-proxy/pkg/jwtAuthenticationMiddleware"
	"github.com/adrichem/go-auth-proxy/pkg/proxy"
	"github.com/adrichem/go-auth-proxy/pkg/tokenextractor"
	"github.com/adrichem/go-auth-proxy/pkg/verifyaudience"
	"github.com/adrichem/go-auth-proxy/pkg/verifyissuer"

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
	*ListenAddress = *ListenAddress
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(CorsConfig.ReadConfig("cors.json")))
	r.Use(azureAdJwtTokenValidation("token"))
	if *Aud != "" {
		log.Println("aud claim must have value " + *Aud)
		r.Use(verifyaudience.Verify(*Aud, claimsSelector.FromGinContext("token")))
	}
	if *Iss != "" {
		log.Println("iss claim must have value " + *Iss)
		r.Use(verifyissuer.Verify(*Iss, claimsSelector.FromGinContext("token")))
	}
	if *Upstream != "" {
		log.Printf("Proxying to %s", *Upstream)
		r.Any("/*path", proxy.Proxy(*Upstream, *HeaderName, *HeaderValue))
	} else {
		log.Println("Running in load test mode. All authenticated requests get HTTP 200 response")
		r.Any("/*path", func(c *gin.Context) { c.Status(http.StatusOK) })
	}

	log.Printf("Listening on %s", *ListenAddress)
	err := r.Run(*ListenAddress)
	if err != nil {
		panic(err)
	}
}

func azureAdJwtTokenValidation(paramName string) gin.HandlerFunc {
	//Azure AD keyset is independent from AD tenant
	jwks, err := keyfunc.Get("https://login.microsoftonline.com/common/discovery/v2.0/keys", keyfunc.Options{RefreshInterval: time.Hour * 24})
	if err != nil {
		panic(err)
	}
	extractToken := func(c *gin.Context) (string, error) { return tokenextractor.ExtractTokenFromHttpRequest(c.Request) }
	pass := func(c *gin.Context, t *jwt.Token) { c.Set("token", t); pass(c) }
	return jwtValidator.Create(jwks.Keyfunc, extractToken, pass, fail)
}
