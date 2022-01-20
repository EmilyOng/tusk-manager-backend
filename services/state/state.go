package services

import (
	"errors"
	"log"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	taskService "github.com/EmilyOng/cvwo/backend/services/task"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func CreateState(payload models.CreateStatePayload) models.CreateStateResponse {
	state := models.State{Name: payload.Name, BoardID: &payload.BoardID, CurrentPosition: payload.CurrentPosition}
	result := db.DB.Create(&state)
	if result.Error != nil {
		log.Println(result.Error)
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
	result := db.DB.Model(&models.State{ID: state.ID}).Save(&state)
	if result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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
	tasks, err := getStateTasks(payload.ID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DeleteStateResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}
		return models.DeleteStateResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	for _, task := range tasks {
		deleteTaskResponse := taskService.DeleteTask(models.DeleteTaskPayload{ID: task.ID})
		if len(deleteTaskResponse.Error) > 0 {
			return models.DeleteStateResponse{
				Response: models.Response{Error: deleteTaskResponse.Error},
			}
		}
	}

	result := db.DB.Delete(&state)
	if result.Error != nil {
		log.Print(result.Error)
		return models.DeleteStateResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.DeleteStateResponse{}
}

func getStateTasks(stateID uint8) ([]models.TaskPrimitive, error) {
	state := models.State{ID: stateID}
	var tasks []models.TaskPrimitive
	err := db.DB.Model(&state).Association("Tasks").Find(&tasks)
	return tasks, err
}
