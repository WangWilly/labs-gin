package dltask

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/WangWilly/labs-gin/pkgs/tasks"
	. "github.com/smartystreets/goconvey/convey"
	gomock "go.uber.org/mock/gomock"
)

////////////////////////////////////////////////////////////////////////////////

func TestCreate(t *testing.T) {
	Convey("Given a DLTask controller", t, func(c C) {
		testInit(t, func(suite *testSuite) {
			Convey("When creating a task with invalid request (missing URL)", func() {
				// Arrange
				// Empty request
				reqBody := CreateRequest{}
				var respBody map[string]interface{}

				// Act
				// Execute request
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/dlTask",
					reqBody,
					&respBody,
					http.StatusBadRequest,
				)

				// Assert response
				So(respBody["error"], ShouldEqual, "Key: 'CreateRequest.Url' Error:Field validation for 'Url' failed on the 'required' tag")
			})

			Convey("When creating a task with valid URL", func(c C) {
				// Arrange
				url := "https://example.com/video.mp4"
				reqBody := CreateRequest{
					Url: url,
				}

				// The TaskManager should be called with a task
				filID := "123e4567-e89b-12d3-a456-426614174000"
				suite.uuidGen.EXPECT().New().Return(filID).Times(1)
				taskID := "123e4567-e89b-12d3-a456-5426614174000"
				suite.uuidGen.EXPECT().New().Return(taskID).Times(1)

				ctx := context.Background()
				suite.taskManager.EXPECT().GetCtx().Return(ctx)
				p, err := filepath.Abs("./public/downloads/" + filID + ".mp4")
				So(err, ShouldBeNil)
				ytdTask := tasks.NewRetribleNamedTaskWithCtx(
					ctx,
					taskID,
					url,
					p,
					suite.controller.cfg.RetryDelay,
					suite.controller.cfg.MaxRetries,
				).WithMaxTimeout(
					suite.controller.cfg.MaxTimeout,
				)
				suite.taskManager.EXPECT().SubmitTask(gomock.Any()).Do(func(task *tasks.DownloadTask) {
					// Assert that the task is of type *tasks.DownloadTask
					c.So(task, ShouldHaveSameTypeAs, ytdTask)
					// Assert that the task is the same as the one we created
					c.So(task.GetID(), ShouldEqual, ytdTask.GetID())
					c.So(task.GetFilePath(), ShouldEqual, ytdTask.GetFilePath())
					c.So(task.GetTargetUrl(), ShouldEqual, ytdTask.GetTargetUrl())
					c.So(task.GetMaxTimeout(), ShouldEqual, ytdTask.GetMaxTimeout())
					c.So(task.GetMaxRetries(), ShouldEqual, ytdTask.GetMaxRetries())
					c.So(task.GetRetries(), ShouldEqual, ytdTask.GetRetries())
					c.So(task.GetProgress(), ShouldEqual, ytdTask.GetProgress())
				})

				// Act
				// Request and response
				respBody := &CreateResponse{}
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodPost,
					"/dlTask",
					reqBody,
					respBody,
					http.StatusCreated,
				)

				// Assert response
				So(respBody.TaskID, ShouldEqual, taskID)
				So(respBody.FileID, ShouldEqual, filID+".mp4")
				So(respBody.Status, ShouldEqual, "task submitted")
			})
		})
	})
}
