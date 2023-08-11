package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	userService "github.com/EmilyOng/cvwo/backend/services/user"
	authUtils "github.com/EmilyOng/cvwo/backend/utils/auth"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	seedUtils "github.com/EmilyOng/cvwo/backend/utils/seed"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAuthToken(ctx *gin.Context) (token string) {
	token_ := strings.Split(ctx.Request.Header.Get("Authorization"), "Bearer ")
	if len(token_) < 2 {
		// Missing authentication token
		ctx.Set(authUtils.UserKey, nil)
		return
	}

	token = strings.Trim(token_[1], " ")
	return
}

func IsAuthenticated(ctx *gin.Context) {
	userInterface, _ := ctx.Get(authUtils.UserKey)
	if userInterface == nil {
		// User is unauthenticated
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			errorUtils.MakeResponseErr(models.UnauthorizedError),
		)
		return
	}

	user := userInterface.(models.AuthUser)
	token := GetAuthToken(ctx)
	// Persists user authentication
	ctx.JSON(http.StatusOK, models.AuthUserResponse{
		User: models.AuthUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		},
	})
}

func Login(ctx *gin.Context) {
	var payload models.LoginPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.TypeMismatch),
		)
		return
	}

	var user models.UserPrimitive
	err = db.DB.Model(&models.User{}).Where("email = ?", payload.Email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User record does not exist
		ctx.AbortWithStatusJSON(
			http.StatusOK,
			errorUtils.MakeResponseErr(models.NotFound),
		)
		return
	}

	err = authUtils.ComparePassword(user.Password, payload.Password)
	if err != nil {
		// User input password does not match
		ctx.AbortWithStatusJSON(
			http.StatusOK,
			errorUtils.MakeResponseErr(models.UnauthorizedError),
		)
		return
	}

	signedToken, err := authUtils.GenerateToken(user)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}
	ctx.JSON(http.StatusOK, models.LoginResponse{
		User: models.AuthUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: signedToken,
		},
	})
}

func SignUp(ctx *gin.Context) {
	var payload models.SignUpPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.TypeMismatch),
		)
		return
	}

	user := models.UserPrimitive{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}
	err = db.DB.Model(&models.User{}).Where("email = ?", payload.Email).First(&user).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// User record already exists
		ctx.AbortWithStatusJSON(
			http.StatusOK,
			errorUtils.MakeResponseErr(models.ConflictError),
		)
		return
	}

	hashedPassword, err := authUtils.HashPassword(payload.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	user.Password = hashedPassword
	user, err = userService.CreateUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	signedToken, err := authUtils.GenerateToken(user)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	// Generate seed data
	err = seedUtils.SeedData(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}
	ctx.JSON(http.StatusOK, models.SignUpResponse{
		User: models.AuthUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: signedToken,
		},
	})
}

func Logout(c *gin.Context) {
	c.Set(authUtils.UserKey, nil)
	c.JSON(http.StatusOK, gin.H{})
}
