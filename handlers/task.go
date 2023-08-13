package handlers

import (
	"fmt"
	"net/http"

	"github.com/EmilyOng/cvwo/backend/models"
	taskService "github.com/EmilyOng/cvwo/backend/services/task"
	errorUtils "github.com/EmilyOng/cvwo/backend/utils/error"

	"github.com/gin-gonic/gin"
)

func CreateTask(ctx *gin.Context) {
	var payload models.CreateTaskPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	createTaskResponse := taskService.CreateTask(payload)
	ctx.JSON(errorUtils.MakeResponseCode(createTaskResponse.Response), createTaskResponse)
}

func UpdateTask(ctx *gin.Context) {
	var payload models.UpdateTaskPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			errorUtils.MakeResponseErr(models.ServerError),
		)
		return
	}

	updateTaskResponse := taskService.UpdateTask(payload)
	ctx.JSON(errorUtils.MakeResponseCode(updateTaskResponse.Response), updateTaskResponse)
}

func DeleteTask(ctx *gin.Context) {
	var taskId uint8
	fmt.Sscan(ctx.Param("task_id"), &taskId)

	deleteTaskResponse := taskService.DeleteTask(models.DeleteTaskPayload{ID: taskId})
	ctx.JSON(errorUtils.MakeResponseCode(deleteTaskResponse.Response), deleteTaskResponse)
}
