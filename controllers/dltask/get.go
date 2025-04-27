package dltask

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetStatus(ctx *gin.Context) {
	taskID := ctx.Param("tid")
	if taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task ID is required"})
		return
	}
	taskStatus, err := c.TaskManager.GetTaskProgress(taskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"status":  taskStatus,
	})
}
