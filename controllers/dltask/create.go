package dltask

import (
	"net/http"

	"github.com/WangWilly/labs-gin/pkgs/tasks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

////////////////////////////////////////////////////////////////////////////////

type CreateRequest struct {
	Url string `json:"url" binding:"required"`
}

func (c *Controller) Create(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileID := uuid.New().String() + ".mp4"
	filePath := c.cfg.DlFolderRoot + "/" + fileID
	ytdlTask := tasks.NewTask(req.Url, filePath)
	c.TaskManager.SubmitTask(ytdlTask)
	taskID := ytdlTask.GetID()

	ctx.JSON(http.StatusCreated, gin.H{
		"task_id": taskID,
		"file_id": fileID,
		"status":  "task submitted",
	})
}
