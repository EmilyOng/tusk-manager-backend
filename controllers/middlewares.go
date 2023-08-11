package controllers

import (
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	authUtils "github.com/EmilyOng/cvwo/backend/utils/auth"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func SetAuthUser(ctx *gin.Context) {
	token := GetAuthToken(ctx)

	claims, err := authUtils.ValidateToken(token)

	if err != nil {
		ctx.Set(authUtils.UserKey, nil)
		return
	}

	user := models.AuthUser{
		ID:    claims.UserID,
		Name:  claims.UserName,
		Email: claims.UserEmail,
		Token: token,
	}

	ctx.Set(authUtils.UserKey, user)
}

func AuthGuard(ctx *gin.Context) {
	userInterface, _ := ctx.Get(authUtils.UserKey)
	if userInterface == nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			errorUtils.MakeResponseErr(models.UnauthorizedError),
		)
	}
}
