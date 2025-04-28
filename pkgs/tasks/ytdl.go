package tasks

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/google/uuid"
)

////////////////////////////////////////////////////////////////////////////////

type DownloadTask struct {
	taskID    string
	targetUrl string
	filePath  string
	progress  int64
	ctx       context.Context
	cancel    context.CancelFunc
}

////////////////////////////////////////////////////////////////////////////////

func NewTask(url string, filepath string) *DownloadTask {
	ctx, cancel := context.WithCancel(context.Background())
	task := &DownloadTask{
		taskID:    uuid.New().String(),
		targetUrl: url,
		filePath:  filepath,
		progress:  0,
		ctx:       ctx,
		cancel:    cancel,
	}
	return task
}

func NewTaskWithCtx(ctx context.Context, url string, filepath string) *DownloadTask {
	ctx, cancel := context.WithCancel(ctx)
	task := &DownloadTask{
		taskID:    uuid.New().String(),
		targetUrl: url,
		filePath:  filepath,
		progress:  0,
		ctx:       ctx,
		cancel:    cancel,
	}
	return task
}

////////////////////////////////////////////////////////////////////////////////

func (t *DownloadTask) GetID() string {
	return t.taskID
}

func (t *DownloadTask) GetProgress() int64 {
	return t.progress
}

func (t *DownloadTask) Execute() {
	t.progress = 30

	if err := exec.CommandContext(
		t.ctx,
		"yt-dlp",
		"-o", t.filePath,
		"-f", "mp4",
		t.targetUrl,
	).Run(); err != nil {
		t.progress = -1
		if t.ctx.Err() == context.Canceled {
			fmt.Printf("Download canceled: %s\n", t.filePath)
		} else {
			fmt.Printf("Error executing command: %v\n", err)
		}
		return
	}
	t.progress = 100

	fmt.Printf("Download complete: %s\n", t.filePath)
}

func (t *DownloadTask) Cancel() {
	fmt.Printf("Canceling download: %s, ", t.filePath)
	fmt.Printf("Canceling task: %s\n", t.taskID)
	t.cancel()
}
