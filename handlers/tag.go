package handlers

import (
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	tagService "github.com/EmilyOng/cvwo/backend/services/tag"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func CreateTag(ctx *gin.Context) {
	var payload models.CreateTagPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	createTagResponse := tagService.CreateTag(payload)
	ctx.JSON(errorUtils.MakeResponseCode(createTagResponse.Response), createTagResponse)
}

func DeleteTag(ctx *gin.Context) {
	deleteTagResponse := tagService.DeleteTag(models.DeleteTagPayload{ID: ctx.Param("tag_id")})
	ctx.JSON(errorUtils.MakeResponseCode(deleteTagResponse.Response), deleteTagResponse)
}

func UpdateTag(ctx *gin.Context) {
	var payload models.UpdateTagPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	updateTagResponse := tagService.UpdateTag(payload)
	ctx.JSON(errorUtils.MakeResponseCode(updateTagResponse.Response), updateTagResponse)
}
