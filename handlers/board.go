package handlers

import (
	"net/http"

	boardService "github.com/EmilyOng/tusk-manager/backend/services/board"
	userService "github.com/EmilyOng/tusk-manager/backend/services/user"
	authUtils "github.com/EmilyOng/tusk-manager/backend/utils/auth"
	"github.com/EmilyOng/tusk-manager/backend/views"

	"github.com/gin-gonic/gin"
)

func GetUserBoards(ctx *gin.Context) {
	userInterface, _ := ctx.Get(authUtils.UserKey)
	if userInterface == nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			views.Response{
				Message: unauthorizedMessage,
				Code:    http.StatusUnauthorized,
			},
		)
		return
	}
	authUserView := userInterface.(views.AuthUserView)

	getUserBoardsResponse := userService.GetUserBoards(views.GetUserBoardsPayload{UserID: authUserView.ID})
	ctx.JSON(getUserBoardsResponse.Code, getUserBoardsResponse)
}

func GetBoardTasks(ctx *gin.Context) {
	getBoardTasksResponse := boardService.GetBoardTasks(views.GetBoardTasksPayload{BoardID: ctx.Param("board_id")})
	ctx.JSON(getBoardTasksResponse.Code, getBoardTasksResponse)
}

func CreateBoard(ctx *gin.Context) {
	var payload views.CreateBoardPayload

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

	createBoardResponse := boardService.CreateBoard(payload)
	ctx.JSON(createBoardResponse.Code, createBoardResponse)
}

func GetBoardTags(ctx *gin.Context) {
	getBoardTagsResponse := boardService.GetBoardTags(views.GetBoardTagsPayload{BoardID: ctx.Param("board_id")})
	ctx.JSON(getBoardTagsResponse.Code, getBoardTagsResponse)
}

func GetBoardMemberProfiles(ctx *gin.Context) {
	getBoardMemberProfilesResponse := boardService.GetBoardMemberProfiles(
		views.GetBoardMemberProfilesPayload{BoardID: ctx.Param("board_id")},
	)
	ctx.JSON(getBoardMemberProfilesResponse.Code, getBoardMemberProfilesResponse)
}

func GetBoard(ctx *gin.Context) {
	getBoardResponse := boardService.GetBoard(views.GetBoardPayload{ID: ctx.Param("board_id")})
	ctx.JSON(getBoardResponse.Code, getBoardResponse)
}

func UpdateBoard(ctx *gin.Context) {
	var payload views.UpdateBoardPayload

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

	updateBoardResponse := boardService.UpdateBoard(payload)
	ctx.JSON(updateBoardResponse.Code, updateBoardResponse)
}

func DeleteBoard(ctx *gin.Context) {
	deleteBoardResponse := boardService.DeleteBoard(views.DeleteBoardPayload{ID: ctx.Param("board_id")})
	ctx.JSON(deleteBoardResponse.Code, deleteBoardResponse)
}

func GetBoardStates(ctx *gin.Context) {
	getBoardStatesResponse := boardService.GetBoardStates(views.GetBoardStatesPayload{BoardID: ctx.Param("board_id")})
	ctx.JSON(getBoardStatesResponse.Code, getBoardStatesResponse)
}
