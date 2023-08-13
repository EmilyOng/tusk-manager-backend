package services

import (
	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
)

func CreateUser(user models.User) (models.User, error) {
	result := db.DB.Model(&models.User{}).Create(&user)
	return user, result.Error
}

func FindUser(email string) (models.User, error) {
	var user models.User
	err := db.DB.Model(&models.User{}).Where("email = ?", email).First(&user).Error

	return user, err
}

func GetUserBoards(userId string) ([]models.BoardPrimitive, error) {
	var boards []models.BoardPrimitive
	var boardIds []string

	err := db.DB.Model(&models.Member{}).Where("user_id = ?", userId).Select("board_id").Find(&boardIds).Error
	if err != nil {
		return boards, err
	}
	err = db.DB.Model(&models.Board{}).Where("id IN ?", boardIds).Find(&boards).Error
	return boards, err
}
