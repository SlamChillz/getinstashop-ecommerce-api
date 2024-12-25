package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/constants"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
	"log"
	"net/http"
	"strings"
)

func AuthMiddy(token *token.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(constants.AuthenticationHeader)
		if len(authHeader) <= len(constants.AuthenticationScheme) {
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
		if strings.ToLower(authHeaderValues[0]) != constants.AuthenticationScheme {
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
		ctx.Set(constants.ContextUserIdKey, user.UserID)
		ctx.Set(constants.ContextUserAdminStatusKey, user.Admin)
		ctx.Next()
	}
}

func AdminMiddy(ctx *gin.Context) {
	_, userIdExist := ctx.Get(constants.ContextUserIdKey)
	userAdminStatus, userAdminStatusExist := ctx.Get(constants.ContextUserAdminStatusKey)
	if !userIdExist || !userAdminStatusExist {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized",
			"error":   gin.H{},
		})
		return
	}
	adminStatus, _ := userAdminStatus.(bool)
	if !adminStatus {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Forbidden",
			"error":   gin.H{},
		})
		return
	}
	ctx.Next()
}
