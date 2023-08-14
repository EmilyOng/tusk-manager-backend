package handlers

import (
	"net/http"

	tagService "github.com/EmilyOng/tusk-manager/backend/services/tag"
	"github.com/EmilyOng/tusk-manager/backend/views"

	"github.com/gin-gonic/gin"
)

func CreateTag(ctx *gin.Context) {
	var payload views.CreateTagPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			views.Response{
				Message: typeMismatchErrorMessage,
				Code:    http.StatusBadRequest,
			},
		)
		return
	}

	createTagResponse := tagService.CreateTag(payload)
	ctx.JSON(createTagResponse.Code, createTagResponse)
}

func DeleteTag(ctx *gin.Context) {
	deleteTagResponse := tagService.DeleteTag(views.DeleteTagPayload{ID: ctx.Param("tag_id")})
	ctx.JSON(deleteTagResponse.Code, deleteTagResponse)
}

func UpdateTag(ctx *gin.Context) {
	var payload views.UpdateTagPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			views.Response{
				Message: typeMismatchErrorMessage,
				Code:    http.StatusBadRequest,
			},
		)
		return
	}

	updateTagResponse := tagService.UpdateTag(payload)
	ctx.JSON(updateTagResponse.Code, updateTagResponse)
}
