package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	"log"
	"net/http"
	"strings"
)

const (
	AuthenticationHeader     = "authorization"
	AuthenticationScheme     = "bearer"
	AuthenticationContextKey = "user"
)

func AuthMiddy(token *token.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AuthenticationHeader)
		if len(authHeader) <= len(AuthenticationScheme) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Missing auth header",
				"error":   gin.H{},
			})
			return
		}
		authHeaderValues := strings.Fields(authHeader)
		if len(authHeaderValues) != 2 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Unrecognized auth header format",
				"error":   gin.H{},
			})
			return
		}
		if strings.ToLower(authHeaderValues[0]) != AuthenticationScheme {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Unsupported auth scheme",
				"error":   gin.H{},
			})
			return
		}
		accessToken := authHeaderValues[1]
		user, err := token.VerifyToken(accessToken)
		if err != nil {
			log.Printf("Error verifying token: %s", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid access token",
				"error":   gin.H{},
			})
			return
		}
		ctx.Set(AuthenticationContextKey, user)
		ctx.Next()
	}
}
