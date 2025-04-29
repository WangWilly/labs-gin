package tasks

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

////////////////////////////////////////////////////////////////////////////////

type DownloadTask struct {
	taskID    string
	targetUrl string
	filePath  string
	progress  int64

	retries      int
	retryDelay   time.Duration
	maxRetries   int
	retryChannel chan struct{}

	ctx        context.Context
	cancel     context.CancelFunc
	maxTimeout time.Duration
}

////////////////////////////////////////////////////////////////////////////////

func NewTask(url string, filepath string) *DownloadTask {
	ctx, cancel := context.WithCancel(context.Background())
	task := &DownloadTask{
		taskID:    uuid.New().String(),
		targetUrl: url,
		filePath:  filepath,
		progress:  0,

		retries:      0,
		retryDelay:   0,
		maxRetries:   0,
		retryChannel: make(chan struct{}, 1),

		ctx:        ctx,
		cancel:     cancel,
		maxTimeout: 0,
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

		retries:      0,
		retryDelay:   0,
		maxRetries:   0,
		retryChannel: make(chan struct{}, 1),

		ctx:        ctx,
		cancel:     cancel,
		maxTimeout: 0,
	}
	return task
}

func NewRetribleTaskWithCtx(
	ctx context.Context,
	url string,
	filepath string,
	retryDelay time.Duration,
	maxRetries int,
) *DownloadTask {
	ctx, cancel := context.WithCancel(ctx)
	task := &DownloadTask{
		taskID:    uuid.New().String(),
		targetUrl: url,
		filePath:  filepath,
		progress:  0,

		retries:      0,
		retryDelay:   retryDelay,
		maxRetries:   maxRetries,
		retryChannel: make(chan struct{}, 1),

		ctx:        ctx,
		cancel:     cancel,
		maxTimeout: 0,
	}
	return task
}

////////////////////////////////////////////////////////////////////////////////

func (t *DownloadTask) WithMaxTimeout(timeout time.Duration) *DownloadTask {
	t.maxTimeout = timeout
	return t
}

////////////////////////////////////////////////////////////////////////////////

func (t *DownloadTask) GetID() string {
	return t.taskID
}

func (t *DownloadTask) GetProgress() int64 {
	return t.progress
}

func (t *DownloadTask) Execute() bool {
	// Setup
	t.progress = 30
	ctx := t.ctx
	if t.maxTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(t.ctx, t.maxTimeout)
		defer cancel()
	}

	// Execute
	if err := exec.CommandContext(
		ctx,
		"yt-dlp",
		"-o", t.filePath,
		"-f", "mp4",
		t.targetUrl,
	).Run(); err != nil {
		t.progress = -1
		if t.ctx.Err() == context.Canceled {
			fmt.Printf("Download canceled: %s\n", t.filePath)
		} else {
			if t.retries < t.maxRetries {
				t.progress = -2
			}
			fmt.Printf("Error executing command: %v\n", err)
		}
		return false
	}
	t.progress = 100

	// Cleanup
	fmt.Printf("Download complete: %s\n", t.filePath)
	return true
}

func (t *DownloadTask) SetRetrySignal() <-chan struct{} {
	go func() {
		if t.retries >= t.maxRetries {
			fmt.Printf("Max retries reached for: %s\n", t.filePath)
			return
		}

		time.Sleep(t.retryDelay)
		t.retries++
		fmt.Printf("Retrying download: %s, attempt: %d\n", t.filePath, t.retries)
		t.retryChannel <- struct{}{}
	}()

	return t.retryChannel
}

func (t *DownloadTask) Cancel() {
	fmt.Printf("Canceling download: %s, ", t.filePath)
	fmt.Printf("Canceling task: %s\n", t.taskID)
	t.cancel()
}
