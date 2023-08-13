package handlers

import (
	"net/http"

	taskService "github.com/EmilyOng/cvwo/backend/services/task"
	"github.com/EmilyOng/cvwo/backend/views"

	"github.com/gin-gonic/gin"
)

func CreateTask(ctx *gin.Context) {
	var payload views.CreateTaskPayload

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

	createTaskResponse := taskService.CreateTask(payload)
	ctx.JSON(createTaskResponse.Code, createTaskResponse)
}

func UpdateTask(ctx *gin.Context) {
	var payload views.UpdateTaskPayload

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

	updateTaskResponse := taskService.UpdateTask(payload)
	ctx.JSON(updateTaskResponse.Code, updateTaskResponse)
}

func DeleteTask(ctx *gin.Context) {
	deleteTaskResponse := taskService.DeleteTask(views.DeleteTaskPayload{ID: ctx.Param("task_id")})
	ctx.JSON(deleteTaskResponse.Code, deleteTaskResponse)
}
