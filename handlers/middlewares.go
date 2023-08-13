package handlers

import (
	"net/http"

	authUtils "github.com/EmilyOng/cvwo/backend/utils/auth"
	"github.com/EmilyOng/cvwo/backend/views"

	"github.com/gin-gonic/gin"
)

const (
	unauthorizedMessage      = "Unauthorized!"
	typeMismatchErrorMessage = "Payload does not match expected type."
)

func SetAuthUser(ctx *gin.Context) {
	token := GetAuthToken(ctx)

	claims, err := authUtils.ValidateToken(token)

	if err != nil {
		ctx.Set(authUtils.UserKey, nil)
		return
	}

	authUserView := views.AuthUserView{
		ID:    claims.UserID,
		Name:  claims.UserName,
		Email: claims.UserEmail,
		Token: token,
	}

	ctx.Set(authUtils.UserKey, authUserView)
}

func AuthGuard(ctx *gin.Context) {
	userInterface, _ := ctx.Get(authUtils.UserKey)
	if userInterface == nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			views.Response{
				Message: unauthorizedMessage,
				Code:    http.StatusUnauthorized,
			},
		)
	}
}
