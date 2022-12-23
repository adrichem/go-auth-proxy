package claimverifier

import (
	"github.com/golang-jwt/jwt/v4"
)

func VerifyClaim[TContext any](tokenSelector func(TContext) *jwt.Token,
	predicate func(jwt.Claims) bool,
	pass func(TContext),
	fail func(TContext)) func(TContext) {
	return func(c TContext) {
		claims := tokenSelector(c).Claims
		if predicate(claims) {
			pass(c)
			return
		}
		fail(c)
	}
}
