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
	unableToCreateTagMessage = "Unable to create tag '%s'."
	unableToUpdateTagMessage = "Unable to update tag (%s)."
	unableToGetTagMessage    = "Unable to retrieve tag (%s)."
	unableToDeleteTagMessage = "Unable to delete tag (%s)."
	tagNotFoundMessage       = "The tag cannot be found (%s)."

	successfullyCreatedTagMessage = "Successfully created tag '%s'!"
	successfullyUpdatedTagMessage = "Successfully updated tag '%s'!"
	successfullyDeletedTagMessage = "Successfully deleted tag '%s'!"
)

func CreateTag(payload views.CreateTagPayload) views.CreateTagResponse {
	tag := models.Tag{Name: payload.Name, Color: payload.Color, BoardID: payload.BoardID}
	err := db.DB.Create(&tag).Error

	if err != nil {
		return views.CreateTagResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToCreateTagMessage, payload.Name),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.CreateTagResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyCreatedTagMessage, tag.Name),
			Code:    http.StatusOK,
		},
		Tag: views.TagMinimalView{ID: tag.ID, Name: tag.Name, Color: tag.Color, BoardID: tag.BoardID},
	}
}

func DeleteTag(payload views.DeleteTagPayload) views.DeleteTagResponse {
	tag := models.Tag{ID: payload.ID}
	err := db.DB.Model(&tag).Preload("Tasks").Find(&tag).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.DeleteTagResponse{
				Response: views.Response{
					Message: fmt.Sprintf(tagNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}

		return views.DeleteTagResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToGetTagMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Remove the association between the tag and its tasks
		err = tx.Model(&tag).Association("Tasks").Delete(&tag.Tasks)
		if err != nil {
			return err
		}
		return tx.Delete(&tag).Error
	})

	if err != nil {
		return views.DeleteTagResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToDeleteTagMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}

	return views.DeleteTagResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyDeletedTagMessage, tag.Name),
			Code:    http.StatusOK,
		},
	}
}

func UpdateTag(payload views.UpdateTagPayload) views.UpdateTagResponse {
	tag := models.Tag{ID: payload.ID, Name: payload.Name, BoardID: payload.BoardID, Color: payload.Color}
	err := db.DB.Save(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return views.UpdateTagResponse{
				Response: views.Response{
					Message: fmt.Sprintf(tagNotFoundMessage, payload.ID),
					Code:    http.StatusUnprocessableEntity,
				},
			}
		}

		return views.UpdateTagResponse{
			Response: views.Response{
				Message: fmt.Sprintf(unableToUpdateTagMessage, payload.ID),
				Code:    http.StatusInternalServerError,
			},
		}
	}
	return views.UpdateTagResponse{
		Response: views.Response{
			Message: fmt.Sprintf(successfullyUpdatedTagMessage, tag.Name),
			Code:    http.StatusOK,
		},
		Tag: views.TagMinimalView{ID: tag.ID, Name: tag.Name, BoardID: tag.BoardID, Color: tag.Color},
	}
}
