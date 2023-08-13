package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Color Color  `gorm:"not null" json:"color" ts_type:"Color"`

	Tasks   []*Task   `gorm:"not null" json:"tasks"`  // Tasks belonging to the board
	Tags    []*Tag    `gorm:"not null" json:"tags"`   // Tags belonging to the board
	States  []*State  `gorm:"not null" json:"states"` // States belonging to the board
	Members []*Member `json:"boardMembers"`           // Members belonging to the board
}

type BoardPrimitive struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color Color  `json:"color" ts_type:"Color"`
}

// Get Board
type GetBoardPayload struct {
	ID string `json:"id"`
}

type GetBoardResponse struct {
	Response
	Board BoardPrimitive `json:"data"`
}

// Create Board
type CreateBoardPayload struct {
	Name   string `json:"name"`
	Color  Color  `json:"color" ts_type:"Color"`
	UserID string `json:"userId"`
}

type CreateBoardResponse struct {
	Response
	Board BoardPrimitive `json:"data"`
}

// Update Board
type UpdateBoardPayload struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Color  Color  `json:"color" ts_type:"Color"`
	UserID string `json:"userId"`
}

type UpdateBoardResponse struct {
	Response
	Board BoardPrimitive `json:"data"`
}

// Get Board Tasks
type GetBoardTasksPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardTasksResponse struct {
	Response
	Tasks []Task `json:"data"`
}

// Get Board Tags
type GetBoardTagsPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardTagsResponse struct {
	Response
	Tags []TagPrimitive `json:"data"`
}

// Get Board Member Profiles
type GetBoardMemberProfilesPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardMemberProfilesResponse struct {
	Response
	MemberProfiles []MemberProfile `json:"data"`
}

// Get Board States
type GetBoardStatesPayload struct {
	BoardID string `json:"boardId"`
}

type GetBoardStatesResponse struct {
	Response
	States []StatePrimitive `json:"data"`
}

// Delete Board
type DeleteBoardPayload struct {
	ID string `json:"id"`
}

type DeleteBoardResponse struct {
	Response
}

func (board *Board) BeforeCreate(tx *gorm.DB) (errr error) {
	// Generates a new UUID
	board.ID = uuid.NewString()
	return
}
