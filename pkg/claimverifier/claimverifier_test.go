package claimverifier

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestClaimsAreNil(t *testing.T) {
	passCalled := false
	failCalled := false
	predicateCalled := false
	token := jwt.Token{}
	fnTokenSelector := func(int) *jwt.Token { return &token }
	fnPredicate := func(jwt.Claims) bool { predicateCalled = true; return true }
	fnPass := func(int) { passCalled = true }
	fnFail := func(int) { failCalled = true }
	fnVerifier := VerifyClaim(fnTokenSelector, fnPredicate, fnPass, fnFail)

	fnVerifier(1)
	assert.True(t, predicateCalled)
	assert.True(t, passCalled)
	assert.False(t, failCalled)
}

func TestClaimsAreEmpty(t *testing.T) {
	passCalled := false
	failCalled := false
	predicateCalled := false
	token := jwt.Token{}
	token.Claims = jwt.RegisteredClaims{}
	fnTokenSelector := func(int) *jwt.Token { return &token }
	fnPredicate := func(jwt.Claims) bool { predicateCalled = true; return true }
	fnPass := func(int) { passCalled = true }
	fnFail := func(int) { failCalled = true }
	fnVerifier := VerifyClaim(fnTokenSelector, fnPredicate, fnPass, fnFail)

	fnVerifier(1)
	assert.True(t, predicateCalled)
	assert.True(t, passCalled)
	assert.False(t, failCalled)
}

func TestPredicateReturnsFalse(t *testing.T) {
	passCalled := false
	failCalled := false
	predicateCalled := false
	token := jwt.Token{}
	token.Claims = jwt.RegisteredClaims{}
	fnTokenSelector := func(int) *jwt.Token { return &token }
	fnPredicate := func(jwt.Claims) bool { predicateCalled = true; return false }
	fnPass := func(int) { passCalled = true }
	fnFail := func(int) { failCalled = true }
	fnVerifier := VerifyClaim(fnTokenSelector, fnPredicate, fnPass, fnFail)

	fnVerifier(1)
	assert.True(t, predicateCalled)
	assert.False(t, passCalled)
	assert.True(t, failCalled)
}
