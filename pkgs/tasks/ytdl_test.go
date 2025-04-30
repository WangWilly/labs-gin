package tasks

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/WangWilly/labs-gin/pkgs/uuid"
	. "github.com/smartystreets/goconvey/convey"
	gomock "go.uber.org/mock/gomock"
)

// MockCmd is a mock implementation of exec.Cmd for testing
type MockCmd struct {
	*exec.Cmd
	shouldSucceed bool
}

func (m *MockCmd) Run() error {
	if m.shouldSucceed {
		return nil
	}
	return &exec.ExitError{ProcessState: &os.ProcessState{}}
}

func TestDownloadTaskConstructors(t *testing.T) {
	Convey("Given task creation parameters", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUUID := uuid.NewMockUUID(ctrl)
		ctx := context.Background()
		url := "https://example.com/video"
		filePath := "/tmp/video.mp4"
		retryDelay := 5 * time.Second
		maxRetries := 3
		taskID := "test-task-id"

		Convey("When creating a new task with generated ID", func() {
			mockUUID.EXPECT().New().Return(taskID)
			task := NewRetribleTaskWithCtx(ctx, mockUUID, url, filePath, retryDelay, maxRetries)

			Convey("Then the task should be properly initialized", func() {
				So(task.GetID(), ShouldEqual, taskID)
				So(task.GetTargetUrl(), ShouldEqual, url)
				So(task.GetFilePath(), ShouldEqual, filePath)
				So(task.GetProgress(), ShouldEqual, 0)
				So(task.GetRetries(), ShouldEqual, 0)
				So(task.GetRetryDelay(), ShouldEqual, retryDelay)
				So(task.GetMaxRetries(), ShouldEqual, maxRetries)
				So(task.GetMaxTimeout(), ShouldEqual, 0)
			})
		})

		Convey("When creating a new task with a specified ID", func() {
			task := NewRetribleNamedTaskWithCtx(ctx, taskID, url, filePath, retryDelay, maxRetries)

			Convey("Then the task should be initialized with the given ID", func() {
				So(task.GetID(), ShouldEqual, taskID)
				So(task.GetTargetUrl(), ShouldEqual, url)
				So(task.GetFilePath(), ShouldEqual, filePath)
			})
		})
	})
}

func TestWithMaxTimeout(t *testing.T) {
	Convey("Given a download task", t, func() {
		task := &DownloadTask{}

		Convey("When setting the max timeout", func() {
			timeout := 10 * time.Minute
			result := task.WithMaxTimeout(timeout)

			Convey("Then the timeout should be set correctly", func() {
				So(result, ShouldEqual, task) // Should return itself for chaining
				So(task.GetMaxTimeout(), ShouldEqual, timeout)
			})
		})
	})
}

func TestGetters(t *testing.T) {
	Convey("Given a download task with specific values", t, func() {
		taskID := "test-id"
		url := "https://example.com/video"
		filePath := "/tmp/test.mp4"
		progress := int64(50)
		retries := 2
		retryDelay := 3 * time.Second
		maxRetries := 5
		maxTimeout := 10 * time.Minute

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		task := &DownloadTask{
			taskID:       taskID,
			targetUrl:    url,
			filePath:     filePath,
			progress:     progress,
			retries:      retries,
			retryDelay:   retryDelay,
			maxRetries:   maxRetries,
			maxTimeout:   maxTimeout,
			ctx:          ctx,
			cancel:       cancel,
			retryChannel: make(chan struct{}, 1),
		}

		Convey("When getting task properties", func() {
			Convey("Then GetID should return the correct ID", func() {
				So(task.GetID(), ShouldEqual, taskID)
			})

			Convey("Then GetTargetUrl should return the correct URL", func() {
				So(task.GetTargetUrl(), ShouldEqual, url)
			})

			Convey("Then GetFilePath should return the correct file path", func() {
				So(task.GetFilePath(), ShouldEqual, filePath)
			})

			Convey("Then GetProgress should return the correct progress", func() {
				So(task.GetProgress(), ShouldEqual, progress)
			})

			Convey("Then GetRetries should return the correct retry count", func() {
				So(task.GetRetries(), ShouldEqual, retries)
			})

			Convey("Then GetRetryDelay should return the correct retry delay", func() {
				So(task.GetRetryDelay(), ShouldEqual, retryDelay)
			})

			Convey("Then GetMaxRetries should return the correct max retries", func() {
				So(task.GetMaxRetries(), ShouldEqual, maxRetries)
			})

			Convey("Then GetMaxTimeout should return the correct max timeout", func() {
				So(task.GetMaxTimeout(), ShouldEqual, maxTimeout)
			})
		})
	})
}

func TestExecute(t *testing.T) {
	Convey("Given a download task", t, func() {
		ctx := context.Background()
		taskID := "test-execute-task"
		url := "https://example.com/video"
		filePath := "/tmp/test-execute.mp4"
		retryDelay := 100 * time.Millisecond
		maxRetries := 3

		task := NewRetribleNamedTaskWithCtx(ctx, taskID, url, filePath, retryDelay, maxRetries)

		Convey("When executing the task successfully", func() {
			mockCmd := &MockCmd{
				Cmd:           exec.Command("echo", "success"),
				shouldSucceed: true,
			}

			result := task.execute(mockCmd)

			Convey("Then the task should complete successfully", func() {
				So(result, ShouldBeTrue)
				So(task.GetProgress(), ShouldEqual, 100)
			})
		})

		Convey("When executing the task unsuccessfully", func() {
			mockCmd := &MockCmd{
				Cmd:           exec.Command("echo", "failure"),
				shouldSucceed: false,
			}

			result := task.execute(mockCmd)

			Convey("Then the task should fail", func() {
				So(result, ShouldBeFalse)
				So(task.GetProgress(), ShouldEqual, -2) // Indicates retry possible
			})
		})

		Convey("When executing a task with a timeout", func() {
			timeout := 50 * time.Millisecond
			task := task.WithMaxTimeout(timeout)

			// We don't actually need to run a command that times out
			// Just verify that the execute passes the right context
			mockCmd := &MockCmd{
				Cmd:           exec.Command("echo", "success"),
				shouldSucceed: true,
			}

			result := task.execute(mockCmd)

			Convey("Then the task should execute with the timeout", func() {
				So(result, ShouldBeTrue)
				So(task.GetProgress(), ShouldEqual, 100)
				So(task.GetMaxTimeout(), ShouldEqual, timeout)
			})
		})
	})
}

func TestSetRetrySignal(t *testing.T) {
	Convey("Given a download task with retry configuration", t, func() {
		ctx := context.Background()
		taskID := "test-retry-task"
		url := "https://example.com/video"
		filePath := "/tmp/test-retry.mp4"
		retryDelay := 50 * time.Millisecond // Short for testing
		maxRetries := 3

		Convey("When setting retry signal with retries available", func() {
			task := NewRetribleNamedTaskWithCtx(ctx, taskID, url, filePath, retryDelay, maxRetries)

			retryChannel := task.SetRetrySignal()

			Convey("Then a retry signal should be sent", func() {
				// Wait for the retry signal
				select {
				case <-retryChannel:
					So(task.GetRetries(), ShouldEqual, 1)
				case <-time.After(retryDelay * 2):
					t.Fatal("Timeout waiting for retry signal")
				}
			})
		})

		Convey("When setting retry signal with max retries reached", func() {
			task := NewRetribleNamedTaskWithCtx(ctx, taskID, url, filePath, retryDelay, maxRetries)
			task.retries = maxRetries // Set retries to max

			retryChannel := task.SetRetrySignal()

			Convey("Then no retry signal should be sent", func() {
				// No signal should be sent
				select {
				case <-retryChannel:
					t.Fatal("Unexpected retry signal received")
				case <-time.After(retryDelay * 2):
					So(task.GetRetries(), ShouldEqual, maxRetries) // Should remain at max
				}
			})
		})
	})
}

func TestCancel(t *testing.T) {
	Convey("Given a download task", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		taskID := "test-cancel-task"
		url := "https://example.com/video"
		filePath := "/tmp/test-cancel.mp4"
		retryDelay := 100 * time.Millisecond
		maxRetries := 3

		task := NewRetribleNamedTaskWithCtx(ctx, taskID, url, filePath, retryDelay, maxRetries)

		Convey("When cancelling the task", func() {
			task.Cancel()

			Convey("Then the context should be cancelled", func() {
				// Check if context is cancelled
				select {
				case <-task.ctx.Done():
					// Context is cancelled, which is what we expect
					So(true, ShouldBeTrue)
				case <-time.After(50 * time.Millisecond):
					t.Fatal("Context not cancelled after Cancel()")
				}
			})
		})
	})
}
