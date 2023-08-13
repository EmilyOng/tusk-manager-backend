package handlers

import (
	"fmt"
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
	var stateID uint8
	fmt.Sscan(ctx.Param("state_id"), &stateID)

	deleteStateResponse := stateService.DeleteState(models.DeleteStatePayload{ID: stateID})
	ctx.JSON(errorUtils.MakeResponseCode(deleteStateResponse.Response), deleteStateResponse)
}
