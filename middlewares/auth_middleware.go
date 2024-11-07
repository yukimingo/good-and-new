package middlewares

import (
	"good-and-new/usecases"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authUsecase usecases.IAuthUsecase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(header, "Bearer") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStrings := strings.TrimPrefix(header, "Bearer ")
		user, err := authUsecase.GetUserFromToken(tokenStrings)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
