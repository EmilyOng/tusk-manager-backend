package controllers

import (
	"log"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	utils "github.com/EmilyOng/cvwo/backend/utils/auth"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func SetAuthUser(c *gin.Context) {
	token := GetAuthToken(c)

	secretKey := utils.GetSecretKey()

	jwtAuth := &utils.JWTAuth{SecretKey: secretKey}
	claims, err := jwtAuth.ValidateToken(token)
	if err != nil {
		log.Println(err)
		c.Set("user", nil)
		return
	}
	user := models.AuthUser{
		ID:    claims.UserID,
		Name:  claims.UserName,
		Email: claims.UserEmail,
		Token: token,
	}

	c.Set("user", user)
}

func AuthGuard(c *gin.Context) {
	userInterface, _ := c.Get("user")
	if userInterface == nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			errorUtils.MakeResponseErr(models.UnauthorizedError),
		)
	}
}
