package dltask

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Cancel(ctx *gin.Context) {
	taskID := ctx.Param("tid")
	if taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task ID is required"})
		return
	}

	currStatus, err := c.TaskManager.GetTaskProgress(taskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if currStatus == 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task already completed"})
		return
	}
	if currStatus == -1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task already cancelled"})
		return
	}

	// Cancel the task
	if err := c.TaskManager.CancelTask(taskID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task_id":              taskID,
		"status_before_cancel": currStatus,
		"status":               "task cancelled",
	})
}
