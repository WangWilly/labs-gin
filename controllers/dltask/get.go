package dltask

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type GetStatusResponse struct {
	TaskID string `json:"task_id"`
	Status int64  `json:"status"`
}

func (c *Controller) GetStatus(ctx *gin.Context) {
	fmt.Println("GetStatus")

	taskID := ctx.Param("tid")
	if taskID == "" {
		fmt.Println("task ID is empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task ID is required"})
		return
	}
	taskStatus, err := c.TaskManager.GetTaskProgress(taskID)
	if err != nil {
		fmt.Printf("failed to get task status: %v\n", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	fmt.Printf("task status: %d\n", taskStatus)
	ctx.JSON(http.StatusOK, GetStatusResponse{
		TaskID: taskID,
		Status: taskStatus,
	})
}
