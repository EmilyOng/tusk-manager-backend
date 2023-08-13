package utils

import (
	"errors"
	"os"
	"time"

	"github.com/EmilyOng/cvwo/backend/models"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Claim struct {
	jwt.StandardClaims
	UserID    string
	UserName  string
	UserEmail string
}

const (
	UserKey string = "user"
)

func GenerateToken(user models.User) (signedToken string, err error) {
	// Token expires in 24 hours
	claims := &Claim{
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
		StandardClaims: jwt.StandardClaims{
			// Express in unix milliseconds
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("AUTH_SECRET_KEY")
	signedToken, err = token.SignedString([]byte(secretKey))
	return
}

func ValidateToken(signedToken string) (claims *Claim, err error) {
	secretKey := os.Getenv("AUTH_SECRET_KEY")
	token, err := jwt.ParseWithClaims(signedToken, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return
	}

	claims, valid := token.Claims.(*Claim)
	if !valid {
		err = errors.New("Cannot parse JWT claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT token has expired")
		return
	}
	return
}

func HashPassword(password string) (hashed string, err error) {
	// Uses a hashing cost of 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return
	}
	hashed = string(bytes)
	return
}

func ComparePassword(userPassword string, passwordInput string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(passwordInput))
	return
}
