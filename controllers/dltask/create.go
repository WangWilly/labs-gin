package dltask

import (
	"net/http"
	"path/filepath"

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
	filePath, err := filepath.Abs(filepath.Join(c.cfg.DlFolderRoot, fileID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get absolute path"})
		return
	}
	// ytdlTask := tasks.NewTaskWithCtx(c.TaskManager.GetCtx(), req.Url, filePath)
	ytdlTask := tasks.NewRetribleTaskWithCtx(
		c.TaskManager.GetCtx(),
		req.Url,
		filePath,
		c.cfg.RetryDelay,
		c.cfg.MaxRetries,
	).WithMaxTimeout(
		c.cfg.MaxTimeout,
	)
	c.TaskManager.SubmitTask(ytdlTask)
	taskID := ytdlTask.GetID()

	ctx.JSON(http.StatusCreated, gin.H{
		"task_id": taskID,
		"file_id": fileID,
		"status":  "task submitted",
	})
}
