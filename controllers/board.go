package controllers

import (
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/db"
	"github.com/EmilyOng/cvwo/backend/models"
	boardService "github.com/EmilyOng/cvwo/backend/services/board"
	userService "github.com/EmilyOng/cvwo/backend/services/user"
	authUtils "github.com/EmilyOng/cvwo/backend/utils/auth"
	commonUtils "github.com/EmilyOng/cvwo/backend/utils/common"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func GetUserBoards(ctx *gin.Context) {
	userInterface, _ := ctx.Get(authUtils.UserKey)
	if userInterface == nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			errorUtils.MakeResponseErr(models.UnauthorizedError),
		)
		return
	}
	authUser := userInterface.(models.AuthUser)

	boards, err := userService.GetUserBoards(authUser.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}
	ctx.JSON(http.StatusOK, models.GetUserBoardsResponse{
		Boards: boards,
	})
}

func GetBoardTasks(c *gin.Context) {
	var boardID uint8
	fmt.Sscan(c.Param("board_id"), &boardID)
	getBoardTasksResponse := boardService.GetBoardTasks(models.GetBoardTasksPayload{BoardID: boardID})
	c.JSON(errorUtils.MakeResponseCode(getBoardTasksResponse.Response), getBoardTasksResponse)
}

func CreateBoard(c *gin.Context) {
	var payload models.CreateBoardPayload

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	var states []*models.State
	for i, state := range commonUtils.GetDefaultStates() {
		states = append(states, &models.State{
			Name:            state,
			CurrentPosition: i,
		})
	}

	res := db.DB.Create(&states)
	if err = res.Error; err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			errorUtils.MakeResponseErr(models.ServerError),
		)
	}

	createBoardResponse := boardService.CreateBoard(payload)
	c.JSON(errorUtils.MakeResponseCode(createBoardResponse.Response), createBoardResponse)
}

func GetBoardTags(c *gin.Context) {
	var boardID uint8
	fmt.Sscan(c.Param("board_id"), &boardID)
	getBoardTagsResponse := boardService.GetBoardTags(models.GetBoardTagsPayload{BoardID: boardID})
	c.JSON(errorUtils.MakeResponseCode(getBoardTagsResponse.Response), getBoardTagsResponse)
}

func GetBoardMemberProfiles(c *gin.Context) {
	var boardID uint8
	fmt.Sscan(c.Param("board_id"), &boardID)
	getBoardMemberProfilesResponse := boardService.GetBoardMemberProfiles(
		models.GetBoardMemberProfilesPayload{BoardID: boardID},
	)
	c.JSON(errorUtils.MakeResponseCode(getBoardMemberProfilesResponse.Response), getBoardMemberProfilesResponse)
}

func GetBoard(c *gin.Context) {
	var boardID uint8
	fmt.Sscan(c.Param("board_id"), &boardID)
	getBoardResponse := boardService.GetBoard(models.GetBoardPayload{ID: boardID})
	c.JSON(errorUtils.MakeResponseCode(getBoardResponse.Response), getBoardResponse)
}

func UpdateBoard(c *gin.Context) {
	var payload models.UpdateBoardPayload

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	updateBoardResponse := boardService.UpdateBoard(payload)
	c.JSON(errorUtils.MakeResponseCode(updateBoardResponse.Response), updateBoardResponse)
}

func DeleteBoard(c *gin.Context) {
	var boardID uint8
	fmt.Sscan(c.Param("board_id"), &boardID)
	deleteBoardResponse := boardService.DeleteBoard(models.DeleteBoardPayload{ID: boardID})
	c.JSON(errorUtils.MakeResponseCode(deleteBoardResponse.Response), deleteBoardResponse)
}

func GetBoardStates(c *gin.Context) {
	var boardID uint8
	fmt.Sscan(c.Param("board_id"), &boardID)
	getBoardStatesResponse := boardService.GetBoardStates(models.GetBoardStatesPayload{BoardID: boardID})
	c.JSON(errorUtils.MakeResponseCode(getBoardStatesResponse.Response), getBoardStatesResponse)
}
