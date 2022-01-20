package controllers

import (
	"errors"
	"log"
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

func GetAuthToken(c *gin.Context) (token string) {
	token_ := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")
	if len(token_) < 2 {
		c.Set("user", nil)
		return
	}
	token = strings.Trim(token_[1], " ")
	return
}

func GenerateJWTToken(c *gin.Context, user models.UserPrimitive) (signedToken string, err error) {
	secretKey := authUtils.GetSecretKey()

	jwtAuth := authUtils.JWTAuth{
		SecretKey: secretKey,
	}

	signedToken, err = jwtAuth.GenerateToken(user)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func IsAuthenticated(c *gin.Context) {
	userInterface, _ := c.Get("user")
	if userInterface == nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			errorUtils.MakeResponseErr(models.UnauthorizedError),
		)
		return
	}
	user := userInterface.(models.AuthUser)
	token := GetAuthToken(c)
	c.JSON(http.StatusOK, models.AuthUserResponse{
		User: models.AuthUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		},
	})
}

func Login(c *gin.Context) {
	var payload models.LoginPayload

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.TypeMismatch),
		)
		return
	}

	var user models.UserPrimitive
	err = db.DB.Model(&models.User{}).Where("email = ?", payload.Email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User record does not exist
		c.AbortWithStatusJSON(
			http.StatusOK,
			errorUtils.MakeResponseErr(models.NotFound),
		)
		return
	}

	err = authUtils.ComparePassword(user.Password, payload.Password)
	if err != nil {
		// User input password does not match
		c.AbortWithStatusJSON(
			http.StatusOK,
			errorUtils.MakeResponseErr(models.UnauthorizedError),
		)
		return
	}

	signedToken, err := GenerateJWTToken(c, user)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}
	c.JSON(http.StatusOK, models.LoginResponse{
		User: models.AuthUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: signedToken,
		},
	})
}

func SignUp(c *gin.Context) {
	var payload models.SignUpPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(
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
		c.AbortWithStatusJSON(
			http.StatusOK,
			errorUtils.MakeResponseErr(models.ConflictError),
		)
		return
	}

	hashedPassword, err := authUtils.HashPassword(payload.Password)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	user.Password = hashedPassword
	user, err = userService.CreateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	signedToken, err := GenerateJWTToken(c, user)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	// Generate seed data
	err = seedUtils.SeedData(&user)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}
	c.JSON(http.StatusOK, models.SignUpResponse{
		User: models.AuthUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: signedToken,
		},
	})
}

func Logout(c *gin.Context) {
	c.Set("user", nil)
	c.JSON(http.StatusOK, gin.H{})
}
