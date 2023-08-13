package handlers

import (
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	boardService "github.com/EmilyOng/cvwo/backend/services/board"
	userService "github.com/EmilyOng/cvwo/backend/services/user"
	authUtils "github.com/EmilyOng/cvwo/backend/utils/auth"
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

func GetBoardTasks(ctx *gin.Context) {
	var boardID uint8
	fmt.Sscan(ctx.Param("board_id"), &boardID)

	getBoardTasksResponse := boardService.GetBoardTasks(models.GetBoardTasksPayload{BoardID: boardID})
	ctx.JSON(errorUtils.MakeResponseCode(getBoardTasksResponse.Response), getBoardTasksResponse)
}

func CreateBoard(ctx *gin.Context) {
	var payload models.CreateBoardPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	createBoardResponse := boardService.CreateBoard(payload)
	ctx.JSON(errorUtils.MakeResponseCode(createBoardResponse.Response), createBoardResponse)
}

func GetBoardTags(ctx *gin.Context) {
	var boardID uint8
	fmt.Sscan(ctx.Param("board_id"), &boardID)

	getBoardTagsResponse := boardService.GetBoardTags(models.GetBoardTagsPayload{BoardID: boardID})
	ctx.JSON(errorUtils.MakeResponseCode(getBoardTagsResponse.Response), getBoardTagsResponse)
}

func GetBoardMemberProfiles(ctx *gin.Context) {
	var boardID uint8
	fmt.Sscan(ctx.Param("board_id"), &boardID)

	getBoardMemberProfilesResponse := boardService.GetBoardMemberProfiles(
		models.GetBoardMemberProfilesPayload{BoardID: boardID},
	)
	ctx.JSON(errorUtils.MakeResponseCode(getBoardMemberProfilesResponse.Response), getBoardMemberProfilesResponse)
}

func GetBoard(ctx *gin.Context) {
	var boardID uint8
	fmt.Sscan(ctx.Param("board_id"), &boardID)

	getBoardResponse := boardService.GetBoard(models.GetBoardPayload{ID: boardID})
	ctx.JSON(errorUtils.MakeResponseCode(getBoardResponse.Response), getBoardResponse)
}

func UpdateBoard(ctx *gin.Context) {
	var payload models.UpdateBoardPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	updateBoardResponse := boardService.UpdateBoard(payload)
	ctx.JSON(errorUtils.MakeResponseCode(updateBoardResponse.Response), updateBoardResponse)
}

func DeleteBoard(ctx *gin.Context) {
	var boardID uint8
	fmt.Sscan(ctx.Param("board_id"), &boardID)

	deleteBoardResponse := boardService.DeleteBoard(models.DeleteBoardPayload{ID: boardID})
	ctx.JSON(errorUtils.MakeResponseCode(deleteBoardResponse.Response), deleteBoardResponse)
}

func GetBoardStates(ctx *gin.Context) {
	var boardID uint8
	fmt.Sscan(ctx.Param("board_id"), &boardID)

	getBoardStatesResponse := boardService.GetBoardStates(models.GetBoardStatesPayload{BoardID: boardID})
	ctx.JSON(errorUtils.MakeResponseCode(getBoardStatesResponse.Response), getBoardStatesResponse)
}
