package handlers

import (
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	memberService "github.com/EmilyOng/cvwo/backend/services/member"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func CreateMember(ctx *gin.Context) {
	var payload models.CreateMemberPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	createMemberResponse := memberService.CreateMember(payload)
	ctx.JSON(errorUtils.MakeResponseCode(createMemberResponse.Response), createMemberResponse)
}

func UpdateMember(ctx *gin.Context) {
	var payload models.UpdateMemberPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	updateMemberResponse := memberService.UpdateMember(payload)
	ctx.JSON(errorUtils.MakeResponseCode(updateMemberResponse.Response), updateMemberResponse)
}

func DeleteMember(ctx *gin.Context) {
	var memberID uint8
	fmt.Sscan(ctx.Param("member_id"), &memberID)

	deleteMemberResponse := memberService.DeleteMember(models.DeleteMemberPayload{ID: memberID})
	ctx.JSON(errorUtils.MakeResponseCode(deleteMemberResponse.Response), deleteMemberResponse)
}
