package main

import (
	"context"
	"net/http"

	"github.com/WangWilly/labs-gin/controllers/dltask"
	"github.com/WangWilly/labs-gin/pkgs/taskmanager"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

////////////////////////////////////////////////////////////////////////////////

type envConfig struct {
	Port           string             `env:"PORT,default=8080"`
	Host           string             `env:"HOST,default=0.0.0.0"`
	TaskManagerCfg taskmanager.Config `nv:",prefix=TASK_MENAGER_"`
	DlTaskCtrlCfg  dltask.Config      `env:",prefix=DL_TASK_CTRL_"`
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	ctx := context.Background()

	// Load environment variables
	cfg := &envConfig{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////

	r := gin.Default()

	////////////////////////////////////////////////////////////////////////////

	r.StaticFile("/favicon.ico", "./public/icon/favicon.ico")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	////////////////////////////////////////////////////////////////////////////

	// Initialize task manager
	taskManager := taskmanager.NewTaskPool(cfg.TaskManagerCfg)
	taskManager.Run()

	////////////////////////////////////////////////////////////////////////////

	// Initialize download task controller
	dlTaskCtrl := dltask.NewController(cfg.DlTaskCtrlCfg, taskManager)
	dlTaskCtrl.RegisterRoutes(r)

	////////////////////////////////////////////////////////////////////////////

	r.Run(cfg.Host + ":" + cfg.Port)
}
