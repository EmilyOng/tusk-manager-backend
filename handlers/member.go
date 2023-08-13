package handlers

import (
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	memberService "github.com/EmilyOng/cvwo/backend/services/member"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func CreateMember(c *gin.Context) {
	var payload models.CreateMemberPayload

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	createMemberResponse := memberService.CreateMember(payload)
	c.JSON(errorUtils.MakeResponseCode(createMemberResponse.Response), createMemberResponse)
}

func UpdateMember(c *gin.Context) {
	var payload models.UpdateMemberPayload

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	updateMemberResponse := memberService.UpdateMember(payload)
	c.JSON(errorUtils.MakeResponseCode(updateMemberResponse.Response), updateMemberResponse)
}

func DeleteMember(c *gin.Context) {
	var memberID uint8
	fmt.Sscan(c.Param("member_id"), &memberID)

	deleteMemberResponse := memberService.DeleteMember(models.DeleteMemberPayload{ID: memberID})
	c.JSON(errorUtils.MakeResponseCode(deleteMemberResponse.Response), deleteMemberResponse)
}
