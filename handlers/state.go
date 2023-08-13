package handlers

import (
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	stateService "github.com/EmilyOng/cvwo/backend/services/state"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func CreateState(ctx *gin.Context) {
	var payload models.CreateStatePayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	createStateResponse := stateService.CreateState(payload)
	ctx.JSON(errorUtils.MakeResponseCode(createStateResponse.Response), createStateResponse)
}

func UpdateState(ctx *gin.Context) {
	var payload models.UpdateStatePayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	updateStateResponse := stateService.UpdateState(payload)
	ctx.JSON(errorUtils.MakeResponseCode(updateStateResponse.Response), updateStateResponse)
}

func DeleteState(ctx *gin.Context) {
	deleteStateResponse := stateService.DeleteState(models.DeleteStatePayload{ID: ctx.Param("state_id")})
	ctx.JSON(errorUtils.MakeResponseCode(deleteStateResponse.Response), deleteStateResponse)
}
