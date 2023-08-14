package views

import "github.com/EmilyOng/tusk-manager/backend/models"

type StateMinimalView struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	CurrentPosition int    `json:"currentPosition"`

	BoardID string `json:"boardId"`
}

type StateFullView = models.State

// Create State
type CreateStatePayload struct {
	Name            string `json:"name"`
	BoardID         string `json:"boardId"`
	CurrentPosition int    `json:"currentPosition"`
}

type CreateStateResponse struct {
	Response
	State StateFullView `json:"data"`
}

// Update State
type UpdateStatePayload struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	BoardID         string `json:"boardId"`
	CurrentPosition int    `json:"currentPosition"`
}

type UpdateStateResponse struct {
	Response
	State StateMinimalView `json:"data"`
}

// Delete State
type DeleteStatePayload struct {
	ID string `json:"id"`
}

type DeleteStateResponse struct {
	Response
}
