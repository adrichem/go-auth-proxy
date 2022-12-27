package verifyissuer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Verify(Expected string, claimsSelector func(c *gin.Context) jwt.MapClaims) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims := claimsSelector(context)
		ok := claims.VerifyIssuer(Expected, Expected != "")
		if !ok {
			context.JSON(http.StatusUnauthorized, "Invalid iss claim")
			context.Abort()
		} else {
			context.Next()
		}
	}
}
