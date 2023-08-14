package services

import (
	"fmt"
	"net/http"

	"github.com/EmilyOng/tusk-manager/backend/db"
	"github.com/EmilyOng/tusk-manager/backend/models"
	"github.com/EmilyOng/tusk-manager/backend/views"
	"gorm.io/gorm"
)

const (
	unableToGetUserBoards = "Unable to get boards for user (%s)."
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

func GetUserBoards(payload views.GetUserBoardsPayload) views.GetUserBoardsResponse {
	var boardsView []views.BoardMinimalView
	var boardIds []string

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Member{}).
			Where("user_id = ?", payload.UserID).
			Select("board_id").
			Find(&boardIds).
			Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Board{}).Where("id IN ?", boardIds).Find(&boardsView).Error
		return err
	})

	if err != nil {
		return views.GetUserBoardsResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetUserBoards, payload.UserID),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.GetUserBoardsResponse{
		Response: views.Response{Code: http.StatusOK},
		Boards:   boardsView,
	}
}
