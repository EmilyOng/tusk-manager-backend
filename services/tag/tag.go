package services

import (
	"errors"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func CreateTag(payload models.CreateTagPayload) models.CreateTagResponse {
	tag := models.Tag{Name: payload.Name, Color: payload.Color, BoardID: &payload.BoardID}
	err := db.DB.Create(&tag).Error

	if err != nil {
		return models.CreateTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.CreateTagResponse{
		Tag: models.TagPrimitive{ID: tag.ID, Name: tag.Name, Color: tag.Color, BoardID: tag.BoardID},
	}
}

func DeleteTag(payload models.DeleteTagPayload) models.DeleteTagResponse {
	tag := models.Tag{ID: payload.ID}
	err := db.DB.Model(&tag).Preload("Tasks").Find(&tag).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DeleteTagResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.DeleteTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	// Remove the association between the tag and its tasks
	err = db.DB.Model(&tag).Association("Tasks").Delete(&tag.Tasks)
	if err != nil {
		return models.DeleteTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	err = db.DB.Delete(&tag).Error
	if err != nil {
		return models.DeleteTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.DeleteTagResponse{}
}

func UpdateTag(payload models.UpdateTagPayload) models.UpdateTagResponse {
	tag := models.Tag{ID: payload.ID, Name: payload.Name, BoardID: &payload.BoardID, Color: payload.Color}
	err := db.DB.Save(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UpdateTagResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.UpdateTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.UpdateTagResponse{
		Tag: models.TagPrimitive{ID: tag.ID, Name: tag.Name, BoardID: tag.BoardID, Color: tag.Color},
	}
}
