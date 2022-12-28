package claimsSelector

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func FromGinContext(paramName string) func(*gin.Context) jwt.MapClaims {
	return func(c *gin.Context) jwt.MapClaims {
		value, found := c.Get(paramName)
		if !found {
			panic("Token not found in context")
		}
		return value.(*jwt.Token).Claims.(jwt.MapClaims)
	}
}
