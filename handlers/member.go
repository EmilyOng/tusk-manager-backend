package handlers

import (
	"net/http"

	memberService "github.com/EmilyOng/cvwo/backend/services/member"
	"github.com/EmilyOng/cvwo/backend/views"

	"github.com/gin-gonic/gin"
)

func CreateMember(ctx *gin.Context) {
	var payload views.CreateMemberPayload

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

	createMemberResponse := memberService.CreateMember(payload)
	ctx.JSON(createMemberResponse.Code, createMemberResponse)
}

func UpdateMember(ctx *gin.Context) {
	var payload views.UpdateMemberPayload

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

	updateMemberResponse := memberService.UpdateMember(payload)
	ctx.JSON(updateMemberResponse.Code, updateMemberResponse)
}

func DeleteMember(ctx *gin.Context) {
	deleteMemberResponse := memberService.DeleteMember(views.DeleteMemberPayload{ID: ctx.Param("member_id")})
	ctx.JSON(deleteMemberResponse.Code, deleteMemberResponse)
}
