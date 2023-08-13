package services

import (
	"errors"
	"time"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	datetime "github.com/EmilyOng/cvwo/backend/utils/datetime"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func CreateTask(payload models.CreateTaskPayload) models.CreateTaskResponse {
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
		return models.CreateTaskResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.CreateTaskResponse{
		Task: task,
	}
}

func getTask(taskId string) (models.Task, error) {
	task := models.Task{ID: taskId}
	result := db.DB.Model(&task).Preload("Tags").Find(&task)
	return task, result.Error
}

func UpdateTask(payload models.UpdateTaskPayload) models.UpdateTaskResponse {
	var tags []*models.Tag
	for _, tag := range payload.Tags {
		tags = append(tags, &models.Tag{
			ID:      tag.ID,
			Name:    tag.Name,
			Color:   tag.Color,
			BoardID: tag.BoardID,
		})
	}

	task, err := getTask(payload.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UpdateTaskResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}
		return models.UpdateTaskResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		task.Name = payload.Name
		task.Description = payload.Description
		task.StateID = &payload.StateID
		task.BoardID = &payload.BoardID
		task.UserID = &payload.UserID

		if len(payload.DueAt) > 0 {
			dueAt, _ := time.Parse(datetime.DatetimeLayout, payload.DueAt)
			task.DueAt = &dueAt
		}
		err := tx.Save(&task).Error
		if err != nil {
			return err
		}

		err = tx.Model(&task).Association("Tags").Replace(&tags)
		return err
	})

	if err != nil {
		return models.UpdateTaskResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.UpdateTaskResponse{
		Task: task,
	}
}

func DeleteTask(payload models.DeleteTaskPayload) models.DeleteTaskResponse {
	task, err := getTask(payload.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DeleteTaskResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}
		return models.DeleteTaskResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	// Remove the association between the task and its tags
	err = db.DB.Model(&task).Association("Tags").Delete(&task.Tags)
	if err != nil {
		return models.DeleteTaskResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	err = db.DB.Delete(&task).Error
	if err != nil {
		return models.DeleteTaskResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.DeleteTaskResponse{}
}
