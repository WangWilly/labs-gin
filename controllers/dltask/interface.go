package dltask

import (
	"context"

	"github.com/WangWilly/labs-gin/pkgs/taskmanager"
)

////////////////////////////////////////////////////////////////////////////////

//go:generate mockgen -source=interface.go -destination=taskmanager_mock.go -package=dltask
type TaskManager interface {
	GetCtx() context.Context
	SubmitTask(taskmanager.Task)
	GetTaskProgress(taskID string) (int64, error)
	CancelTask(taskID string) error
}
