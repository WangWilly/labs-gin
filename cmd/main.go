package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WangWilly/labs-gin/controllers/dltask"
	"github.com/WangWilly/labs-gin/pkgs/taskmanager"
	"github.com/WangWilly/labs-gin/pkgs/utils"
	"github.com/WangWilly/labs-gin/pkgs/uuid"

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

	// Initialize Gin router
	r := utils.GetDefaultRouter()

	////////////////////////////////////////////////////////////////////////////

	// Initialize task manager
	taskManager := taskmanager.NewTaskPool(cfg.TaskManagerCfg)
	taskManager.Run()

	// UUid generator
	uuidGen := uuid.NewGenerator()

	////////////////////////////////////////////////////////////////////////////

	// Initialize download task controller
	dlTaskCtrl := dltask.NewController(cfg.DlTaskCtrlCfg, taskManager, uuidGen)
	dlTaskCtrl.RegisterRoutes(r)

	////////////////////////////////////////////////////////////////////////////

	// Set up the server
	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	////////////////////////////////////////////////////////////////////////////

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// Kill (no param) default sends syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Log shutdown message
	fmt.Println("Received shutdown signal, shutting down server...")
	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown task manager
	fmt.Println("Shutdown Task Manager ...")
	taskManager.ShutdownNow()

	// Gracefully shutdown the server
	fmt.Println("Shutdown HTTP Server ...")
	if err := srv.Shutdown(ctx); err != nil {
		// Handle shutdown error
		panic(err)
	}

	// Wait for tasks to finish or timeout
	<-ctx.Done()
	fmt.Println("Server shutdown complete.")
}
