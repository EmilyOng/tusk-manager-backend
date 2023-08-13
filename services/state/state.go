package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	"github.com/EmilyOng/cvwo/backend/views"
	"gorm.io/gorm"
)

const (
	unableToCreateStateMessage = "Unable to create state '%s'."
	unableToUpdateStateMessage = "Unable to update state (%s)."
	unableToGetStateMessage    = "Unable to retrieve state (%s)."
	unableToDeleteStateMessage = "Unable to delete state (%s)."
	stateNotFoundMessage       = "The state cannot be found (%s)."

	successfullyCreatedStateMessage = "Successfully created state '%s'!"
	successfullyUpdatedStateMessage = "Successfully updated state '%s'!"
	successfullyDeletedStateMessage = "Successfully deleted state '%s'!"
)

func CreateState(payload views.CreateStatePayload) views.CreateStateResponse {
	state := models.State{Name: payload.Name, BoardID: payload.BoardID, CurrentPosition: payload.CurrentPosition}
	err := db.DB.Create(&state).Error

	if err != nil {
		return views.CreateStateResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToCreateStateMessage, payload.Name),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.CreateStateResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyCreatedStateMessage, state.Name),
			Code:    http.StatusOK,
		},
		State: state,
	}
}

func UpdateState(payload views.UpdateStatePayload) views.UpdateStateResponse {
	state := models.State{
		ID:              payload.ID,
		Name:            payload.Name,
		BoardID:         payload.BoardID,
		CurrentPosition: payload.CurrentPosition,
	}
	err := db.DB.Save(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.UpdateStateResponse{
				Response: views.Response{
					Message: fmt.Sprintf(stateNotFoundMessage, payload.ID),
					Code:    http.StatusInternalServerError,
				},
			}
		}

		return views.UpdateStateResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToUpdateStateMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.UpdateStateResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyUpdatedStateMessage, state.Name),
			Code:    http.StatusOK,
		},
		State: views.StateMinimalView{
			ID:              state.ID,
			Name:            state.Name,
			BoardID:         state.BoardID,
			CurrentPosition: state.CurrentPosition,
		},
	}
}

func DeleteState(payload views.DeleteStatePayload) views.DeleteStateResponse {
	state := models.State{ID: payload.ID}

	// Get tasks associated with the state
	var tasksView []views.TaskMinimalView
	err := db.DB.Model(&state).Association("Tasks").Find(&tasksView)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.DeleteStateResponse{
				Response: views.Response{
					Message: fmt.Sprintf(stateNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}
		return views.DeleteStateResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetStateMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// TODO: Delegate tasks to a default state
		return tx.Delete(&state).Error
	})

	if err != nil {
		return views.DeleteStateResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToDeleteStateMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.DeleteStateResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyDeletedStateMessage, state.Name),
			Code:    http.StatusOK,
		},
	}
}
