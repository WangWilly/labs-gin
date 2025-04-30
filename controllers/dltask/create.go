package dltask

import (
	"net/http"
	"path/filepath"

	"github.com/WangWilly/labs-gin/pkgs/tasks"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type CreateRequest struct {
	Url string `json:"url" binding:"required"`
}

type CreateResponse struct {
	TaskID string `json:"task_id"`
	FileID string `json:"file_id"`
	Status string `json:"status"`
}

func (c *Controller) Create(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileID := c.UuidGen.New() + ".mp4"
	filePath, err := filepath.Abs(filepath.Join(c.cfg.DlFolderRoot, fileID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get absolute path"})
		return
	}
	// ytdlTask := tasks.NewTaskWithCtx(c.TaskManager.GetCtx(), req.Url, filePath)
	ytdlTask := tasks.NewRetribleTaskWithCtx(
		c.TaskManager.GetCtx(),
		c.UuidGen,
		req.Url,
		filePath,
		c.cfg.RetryDelay,
		c.cfg.MaxRetries,
	).WithMaxTimeout(
		c.cfg.MaxTimeout,
	)
	c.TaskManager.SubmitTask(ytdlTask)
	taskID := ytdlTask.GetID()

	ctx.JSON(http.StatusCreated, CreateResponse{
		TaskID: taskID,
		FileID: fileID,
		Status: "task submitted",
	})
}
