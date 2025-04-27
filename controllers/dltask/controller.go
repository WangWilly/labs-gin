package dltask

import (
	"github.com/WangWilly/labs-gin/pkgs/taskmanager"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type Config struct {
	DlFolderRoot string `env:"DL_FOLDER_ROOT,default=./public/downloads"`
}

type Controller struct {
	cfg         Config
	TaskManager *taskmanager.TaskPool
}

func NewController(cfg Config, taskManager *taskmanager.TaskPool) *Controller {
	return &Controller{
		cfg:         cfg,
		TaskManager: taskManager,
	}
}

func (c *Controller) RegisterRoutes(r *gin.Engine) {
	////////////////////////////////////////////////////////////////////////////
	// File download
	r.GET("/dlTaskFile/:fid", c.GetFile)

	////////////////////////////////////////////////////////////////////////////
	// Task management
	r.POST("/dlTask", c.Create)
	r.GET("/dlTask/:tid", c.GetStatus)
	r.DELETE("/dlTask/:tid", c.Cancel)
}
