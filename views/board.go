package views

import (
	colorTypes "github.com/EmilyOng/cvwo/backend/types/color"

	"github.com/EmilyOng/cvwo/backend/models"
)

type BoardMinimalView struct {
	ID    string           `json:"id"`
	Name  string           `json:"name"`
	Color colorTypes.Color `json:"color" ts_type:"Color"`
}

type BoardFullView = models.Board

// Get Board
type GetBoardPayload struct {
	ID string `json:"id"`
}

type GetBoardResponse struct {
	Response
	Board BoardMinimalView `json:"data"`
}

// Create Board
type CreateBoardPayload struct {
	Name   string           `json:"name"`
	Color  colorTypes.Color `json:"color" ts_type:"Color"`
	UserID string           `json:"userId"`
}

type CreateBoardResponse struct {
	Response
	Board BoardMinimalView `json:"data"`
}

// Update Board
type UpdateBoardPayload struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Color  colorTypes.Color `json:"color" ts_type:"Color"`
	UserID string           `json:"userId"`
}

type UpdateBoardResponse struct {
	Response
	Board BoardMinimalView `json:"data"`
}

// Get Board Tasks
type GetBoardTasksPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardTasksResponse struct {
	Response
	Tasks []TaskFullView `json:"data"`
}

// Get Board Tags
type GetBoardTagsPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardTagsResponse struct {
	Response
	Tags []TagMinimalView `json:"data"`
}

// Get Board Member Profiles
type GetBoardMemberProfilesPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardMemberProfilesResponse struct {
	Response
	Members []MemberFullView `json:"data"`
}

// Get Board States
type GetBoardStatesPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardStatesResponse struct {
	Response
	States []StateMinimalView `json:"data"`
}

// Delete Board
type DeleteBoardPayload struct {
	ID string `json:"id"`
}

type DeleteBoardResponse struct {
	Response
}
