package jwtauthenticationmiddleware

import (
	"github.com/golang-jwt/jwt/v4"
)

func Create[TContext any](selectKeySet jwt.Keyfunc,
	extractToken func(TContext) (string, error),
	pass func(TContext, *jwt.Token),
	fail func(TContext, error)) func(TContext) {
	return func(c TContext) {
		var parsedToken *jwt.Token
		token, err := extractToken(c)
		if err == nil {
			parsedToken, err = jwt.Parse(token, selectKeySet)
		}
		if err != nil {
			fail(c, err)
			return
		}
		pass(c, parsedToken)
	}
}
