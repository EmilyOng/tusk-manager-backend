package views

import "github.com/EmilyOng/tusk-manager/backend/models"

type GetUserBoardsPayload struct {
	UserID string `json:"userId"`
}

type GetUserBoardsResponse struct {
	Response
	Boards []BoardMinimalView `json:"data"`
}

type UserMinimalView struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserView = models.User
