package handlers

import (
	"net/http"

	stateService "github.com/EmilyOng/cvwo/backend/services/state"
	"github.com/EmilyOng/cvwo/backend/views"

	"github.com/gin-gonic/gin"
)

func CreateState(ctx *gin.Context) {
	var payload views.CreateStatePayload

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

	createStateResponse := stateService.CreateState(payload)
	ctx.JSON(createStateResponse.Code, createStateResponse)
}

func UpdateState(ctx *gin.Context) {
	var payload views.UpdateStatePayload

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

	updateStateResponse := stateService.UpdateState(payload)
	ctx.JSON(updateStateResponse.Code, updateStateResponse)
}

func DeleteState(ctx *gin.Context) {
	deleteStateResponse := stateService.DeleteState(views.DeleteStatePayload{ID: ctx.Param("state_id")})
	ctx.JSON(deleteStateResponse.Code, deleteStateResponse)
}
