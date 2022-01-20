package services

import (
	"errors"
	"log"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"
	"gorm.io/gorm"
)

func CreateTag(payload models.CreateTagPayload) models.CreateTagResponse {
	tag := models.Tag{Name: payload.Name, Color: payload.Color, BoardID: &payload.BoardID}
	result := db.DB.Create(&tag)
	if result.Error != nil {
		log.Println(result.Error)
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
	result := db.DB.Model(&tag).Preload("Tasks").Find(&tag)
	if result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.DeleteTagResponse{
				Response: errorUtils.MakeResponseErr(models.NotFound),
			}
		}

		return models.DeleteTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	err := db.DB.Model(&tag).Association("Tasks").Delete(&tag.Tasks)
	if err != nil {
		log.Println(err)
		return models.DeleteTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	result = db.DB.Delete(&tag)
	if result.Error != nil {
		log.Print(result.Error)
		return models.DeleteTagResponse{
			Response: errorUtils.MakeResponseErr(models.ServerError),
		}
	}
	return models.DeleteTagResponse{}
}

func UpdateTag(payload models.UpdateTagPayload) models.UpdateTagResponse {
	tag := models.Tag{ID: payload.ID, Name: payload.Name, BoardID: &payload.BoardID, Color: payload.Color}
	result := db.DB.Model(&models.Tag{ID: tag.ID}).Save(&tag)
	if result.Error != nil {
		log.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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
