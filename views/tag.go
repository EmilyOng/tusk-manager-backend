package views

import (
	colorTypes "github.com/EmilyOng/cvwo/backend/types/color"
)

type TagMinimalView struct {
	ID    string           `json:"id"`
	Name  string           `json:"name"`
	Color colorTypes.Color `json:"color" ts_type:"Color"`

	BoardID string `json:"boardId"`
}

// Create Tag
type CreateTagPayload struct {
	Name    string           `json:"name"`
	Color   colorTypes.Color `json:"color" ts_type:"Color"`
	BoardID string           `json:"boardId"`
}

type CreateTagResponse struct {
	Response
	Tag TagMinimalView `json:"data"`
}

// Update Tag
type UpdateTagPayload struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	BoardID string           `json:"boardId"`
	Color   colorTypes.Color `json:"color" ts_type:"Color"`
}

type UpdateTagResponse struct {
	Response
	Tag TagMinimalView `json:"data"`
}

// Delete Tag
type DeleteTagPayload struct {
	ID string `json:"id"`
}

type DeleteTagResponse struct {
	Response
}
