package jwtauthenticationmiddleware

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var passCalled bool = false
var failCalled bool = false
var lastError error = nil
var lastToken *jwt.Token = nil
var lastContextValue = 0
var pass = func(c int, t *jwt.Token) { passCalled = true; lastToken = t; lastContextValue = c }
var fail = func(c int, e error) { failCalled = true; lastError = e; lastContextValue = c }

func initDefaults() {
	passCalled = false
	failCalled = false
	lastError = nil
	lastToken = nil
	lastContextValue = 0
}

func defaultClaims(lifeTime int64) jwt.MapClaims {
	claims := jwt.MapClaims{}
	claims["nbf"] = time.Now().Unix()
	claims["iat"] = claims["nbf"]
	claims["exp"] = claims["nbf"].(int64) + lifeTime
	return claims
}

func getSignedString(token *jwt.Token, secret interface{}) string {
	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(tokenString)
	}
	return tokenString
}

func selectorFor[TContext any](value string) func(TContext) (string, error) {
	return func(TContext) (string, error) { return value, nil }
}

func keyFuncFor(key interface{}) func(*jwt.Token) (interface{}, error) {
	return func(*jwt.Token) (interface{}, error) { return key, nil }
}

func TestMultipleContexts(t *testing.T) {
	initDefaults()
	var hmacSampleSecret []byte
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, defaultClaims(10))
	tokenString := getSignedString(token, hmacSampleSecret)
	fnValidate := Create(keyFuncFor(hmacSampleSecret), selectorFor[int](tokenString), pass, fail)
	fnValidate(1)
	assert.False(t, failCalled)
	assert.True(t, passCalled)
	assert.Equal(t, 1, lastContextValue)

	initDefaults()
	fnValidate(2)
	assert.False(t, failCalled)
	assert.True(t, passCalled)
	assert.Equal(t, 2, lastContextValue)

	initDefaults()
	fnValidate(1)
	assert.False(t, failCalled)
	assert.True(t, passCalled)
	assert.Equal(t, 1, lastContextValue)
}

func TestClaimNbf(t *testing.T) {
	initDefaults()
	var hmacSampleSecret []byte
	claims := defaultClaims(10)
	claims["nbf"] = time.Now().AddDate(1, 1, 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString := getSignedString(token, hmacSampleSecret)
	fnValidate := Create(keyFuncFor(hmacSampleSecret), selectorFor[int](tokenString), pass, fail)
	fnValidate(1)
	assert.True(t, failCalled)
	assert.False(t, passCalled)
}

func TestClaimExp(t *testing.T) {
	initDefaults()
	var hmacSampleSecret []byte
	claims := defaultClaims(-100)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString := getSignedString(token, hmacSampleSecret)
	fnValidate := Create(keyFuncFor(hmacSampleSecret), selectorFor[int](tokenString), pass, fail)
	fnValidate(1)
	assert.True(t, failCalled)
	assert.False(t, passCalled)
}

func TestTokenValidSyntax(t *testing.T) {
	initDefaults()
	var hmacSampleSecret []byte
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, defaultClaims(10))
	tokenString := getSignedString(token, hmacSampleSecret)
	fnValidate := Create(keyFuncFor(hmacSampleSecret), selectorFor[int](tokenString), pass, fail)
	fnValidate(1)
	assert.False(t, failCalled)
	assert.True(t, passCalled)
	assert.Equal(t, 1, lastContextValue)
}

func TestTokenInvalidSyntax(t *testing.T) {
	initDefaults()
	var hmacSampleSecret []byte
	fnValidate := Create(keyFuncFor(hmacSampleSecret), selectorFor[int]("INVALID TOKEN STRING"), pass, fail)
	fnValidate(1)
	assert.True(t, failCalled)
	assert.False(t, passCalled)
	assert.Equal(t, 1, lastContextValue)
}

func TestTokenError(t *testing.T) {
	initDefaults()
	var hmacSampleSecret []byte
	fnValidate := Create(keyFuncFor(hmacSampleSecret), func(int) (string, error) { return "", errors.New("XXX") }, pass, fail)
	fnValidate(1)
	assert.True(t, failCalled)
	assert.False(t, passCalled)
	assert.Equal(t, "XXX", lastError.Error())
}
