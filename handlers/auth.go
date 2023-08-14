package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/EmilyOng/tusk-manager/backend/models"
	userService "github.com/EmilyOng/tusk-manager/backend/services/user"
	authUtils "github.com/EmilyOng/tusk-manager/backend/utils/auth"
	seedUtils "github.com/EmilyOng/tusk-manager/backend/utils/seed"
	"github.com/EmilyOng/tusk-manager/backend/views"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	userNotFoundMessage                        = "The email '%s' does not exist. Have you created an account?"
	userAlreadyExistsMessage                   = "The user with the email '%s' already exists."
	unableToCreateUserMessage                  = "Unable to create user '%s'."
	passwordsMismatchErrorMessage              = "Passwords do not match, please try again."
	unableToHashPasswordMessage                = "Unable to hash the password."
	unableToGenerateAuthenticationTokenMessage = "Unable to generate authentication token."
	unableToGenerateSeedDataMessage            = "Unable to generate seed data."

	successfullyLoginMessage  = "Welcome back %s!"
	successfullySignUpMessage = "Welcome %s! Setting things up..."
	successfullyLogoutMessage = "Goodbye!"
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
			views.Response{
				Message: unauthorizedMessage,
				Code:    http.StatusUnauthorized,
			},
		)
		return
	}

	authUserView := userInterface.(views.AuthUserView)
	authUserView.Token = GetAuthToken(ctx)
	// Persists user authentication
	ctx.JSON(http.StatusOK, views.AuthUserResponse{
		Response: views.Response{
			Code: http.StatusOK,
		},
		User: authUserView,
	})
}

func Login(ctx *gin.Context) {
	var payload views.LoginPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			views.Response{
				Message: typeMismatchErrorMessage,
				Code:    http.StatusBadRequest,
			},
		)
		return
	}

	user, err := userService.FindUser(payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User record does not exist
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			views.Response{
				Message: fmt.Sprintf(userNotFoundMessage, payload.Email),
				Code:    http.StatusUnprocessableEntity,
			},
		)
		return
	}

	err = authUtils.ComparePassword(user.Password, payload.Password)
	if err != nil {
		// User input password does not match
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			views.Response{
				Message: passwordsMismatchErrorMessage,
				Code:    http.StatusUnauthorized,
			},
		)
		return
	}

	signedToken, err := authUtils.GenerateToken(user)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			views.Response{
				Message: unableToGenerateAuthenticationTokenMessage,
				Code:    http.StatusInternalServerError,
			},
		)
		return
	}
	ctx.JSON(http.StatusOK, views.LoginResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyLoginMessage, user.Name),
			Code:    http.StatusOK,
		},
		User: views.AuthUserView{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: signedToken,
		},
	})
}

func SignUp(ctx *gin.Context) {
	var payload views.SignUpPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			views.Response{
				Message: typeMismatchErrorMessage,
				Code:    http.StatusBadRequest,
			},
		)
		return
	}

	_, err = userService.FindUser(payload.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// User record already exists
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			views.Response{
				Message: fmt.Sprintf(userAlreadyExistsMessage, payload.Email),
				Code:    http.StatusUnprocessableEntity,
			},
		)
		return
	}

	hashedPassword, err := authUtils.HashPassword(payload.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			views.Response{
				Message: unableToHashPasswordMessage,
				Code:    http.StatusInternalServerError,
			},
		)
		return
	}

	user, err := userService.CreateUser(models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			views.Response{
				Message: fmt.Sprintf(unableToCreateUserMessage, payload.Email),
				Code:    http.StatusInternalServerError,
			},
		)
		return
	}

	signedToken, err := authUtils.GenerateToken(user)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			views.Response{
				Message: unableToGenerateAuthenticationTokenMessage,
				Code:    http.StatusInternalServerError,
			},
		)
		return
	}

	// Generate seed data
	err = seedUtils.SeedData(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			views.Response{
				Message: unableToGenerateSeedDataMessage,
				Code:    http.StatusInternalServerError,
			},
		)
		return
	}
	ctx.JSON(http.StatusOK, views.SignUpResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullySignUpMessage, user.Name),
			Code:    http.StatusOK,
		},
		User: views.AuthUserView{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: signedToken,
		},
	})
}

func Logout(c *gin.Context) {
	c.Set(authUtils.UserKey, nil)
	c.JSON(http.StatusOK, views.Response{
		Message: successfullyLogoutMessage,
		Code:    http.StatusOK,
	})
}
