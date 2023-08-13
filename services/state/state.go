package services

import (
	"errors"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	taskService "github.com/EmilyOng/cvwo/backend/services/task"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func CreateState(payload models.CreateStatePayload) models.CreateStateResponse {
	state := models.State{Name: payload.Name, BoardID: &payload.BoardID, CurrentPosition: payload.CurrentPosition}
	err := db.DB.Create(&state).Error

	if err != nil {
		return models.CreateStateResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.CreateStateResponse{
		State: state,
	}
}

func UpdateState(payload models.UpdateStatePayload) models.UpdateStateResponse {
	state := models.State{
		ID:              payload.ID,
		Name:            payload.Name,
		BoardID:         &payload.BoardID,
		CurrentPosition: payload.CurrentPosition,
	}
	err := db.DB.Save(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UpdateStateResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.UpdateStateResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.UpdateStateResponse{
		State: models.StatePrimitive{
			ID:              state.ID,
			Name:            state.Name,
			BoardID:         state.BoardID,
			CurrentPosition: state.CurrentPosition,
		},
	}
}

func DeleteState(payload models.DeleteStatePayload) models.DeleteStateResponse {
	state := models.State{ID: payload.ID}

	// Get tasks associated with the state
	var tasks []models.TaskPrimitive
	err := db.DB.Model(&state).Association("Tasks").Find(&tasks)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DeleteStateResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}
		return models.DeleteStateResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	// TODO: Delegate to a default state
	for _, task := range tasks {
		deleteTaskResponse := taskService.DeleteTask(models.DeleteTaskPayload{ID: task.ID})
		if len(deleteTaskResponse.Error) > 0 {
			return models.DeleteStateResponse{
				Response: models.Response{Error: deleteTaskResponse.Error},
			}
		}
	}

	err = db.DB.Delete(&state).Error
	if err != nil {
		return models.DeleteStateResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.DeleteStateResponse{}
}
