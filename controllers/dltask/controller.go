package dltask

import (
	"time"

	"github.com/WangWilly/labs-gin/pkgs/uuid"
	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////

type Config struct {
	DlFolderRoot string        `env:"DL_FOLDER_ROOT,default=./public/downloads"`
	RetryDelay   time.Duration `env:"RETRY_DELAY,default=5s"`
	MaxRetries   int           `env:"MAX_RETRIES,default=3"`
	MaxTimeout   time.Duration `env:"MAX_TIMEOUT,default=5m"`
}

type Controller struct {
	cfg         Config
	TaskManager TaskManager
	UuidGen     uuid.UUID
}

func NewController(cfg Config, taskManager TaskManager, uuidGen uuid.UUID) *Controller {
	return &Controller{
		cfg:         cfg,
		TaskManager: taskManager,
		UuidGen:     uuidGen,
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
