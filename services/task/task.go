package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/EmilyOng/tusk-manager/backend/db"
	"github.com/EmilyOng/tusk-manager/backend/models"
	datetime "github.com/EmilyOng/tusk-manager/backend/utils/datetime"
	"github.com/EmilyOng/tusk-manager/backend/views"
	"gorm.io/gorm"
)

const (
	unableToCreateTaskMessage = "Unable to create task '%s'."
	unableToUpdateTaskMessage = "Unable to update task (%s)."
	unableToGetTaskMessage    = "Unable to retrieve task (%s)."
	unableToDeleteTaskMessage = "Unable to delete task (%s)."
	taskNotFoundMessage       = "The task cannot be found (%s)."

	successfullyCreatedTaskMessage = "Successfully created task '%s'!"
	successfullyUpdatedTaskMessage = "Successfully updated task '%s'!"
	successfullyDeletedTaskMessage = "Successfully deleted task '%s'!"
)

func getTask(taskId string) (models.Task, error) {
	task := models.Task{ID: taskId}
	result := db.DB.Preload("Tags").Model(&task).Find(&task)
	return task, result.Error
}

func CreateTask(payload views.CreateTaskPayload) views.CreateTaskResponse {
	var tags []*models.Tag
	for _, tag := range payload.Tags {
		tags = append(tags, &models.Tag{
			ID:      tag.ID,
			Name:    tag.Name,
			Color:   tag.Color,
			BoardID: tag.BoardID,
		})
	}
	task := models.Task{
		Name:        payload.Name,
		Description: payload.Description,
		StateID:     payload.StateID,
		Tags:        tags,
		BoardID:     payload.BoardID,
		UserID:      payload.UserID,
	}
	if len(payload.DueAt) > 0 {
		dueAt, _ := time.Parse(datetime.DatetimeLayout, payload.DueAt)
		task.DueAt = &dueAt
	}
	err := db.DB.Create(&task).Error
	if err != nil {
		return views.CreateTaskResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToCreateTaskMessage, payload.Name),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.CreateTaskResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyCreatedTaskMessage, task.Name),
			Code:    http.StatusOK,
		},
		Task: task,
	}
}

func UpdateTask(payload views.UpdateTaskPayload) views.UpdateTaskResponse {
	task, err := getTask(payload.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.UpdateTaskResponse{
				Response: views.Response{
					Message: fmt.Sprintf(taskNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}
		return views.UpdateTaskResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetTaskMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var tags []*models.Tag
		for _, tag := range payload.Tags {
			tags = append(tags, &models.Tag{
				ID:      tag.ID,
				Name:    tag.Name,
				Color:   tag.Color,
				BoardID: tag.BoardID,
			})
		}

		task.Name = payload.Name
		task.Description = payload.Description
		task.StateID = payload.StateID
		task.BoardID = payload.BoardID
		task.UserID = payload.UserID

		if len(payload.DueAt) > 0 {
			dueAt, _ := time.Parse(datetime.DatetimeLayout, payload.DueAt)
			task.DueAt = &dueAt
		}

		err := tx.Model(&models.Task{ID: task.ID}).Save(&task).Error
		if err != nil {
			return err
		}

		err = tx.Model(&task).Association("Tags").Replace(&tags)
		return err
	})

	if err != nil {
		return views.UpdateTaskResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToUpdateTaskMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.UpdateTaskResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyUpdatedTaskMessage, task.Name),
			Code:    http.StatusOK,
		},
		Task: task,
	}
}

func DeleteTask(payload views.DeleteTaskPayload) views.DeleteTaskResponse {
	task, err := getTask(payload.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.DeleteTaskResponse{
				Response: views.Response{
					Message: fmt.Sprintf(taskNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}
		return views.DeleteTaskResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetTaskMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Remove the association between the task and its tags
		err := tx.Model(&task).Association("Tags").Delete(&task.Tags)
		if err != nil {
			return err
		}
		return tx.Delete(&task).Error
	})

	if err != nil {
		return views.DeleteTaskResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToDeleteTaskMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.DeleteTaskResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyDeletedTaskMessage, task.Name),
			Code:    http.StatusOK,
		},
	}
}
