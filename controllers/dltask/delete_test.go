package dltask

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

func TestCancel(t *testing.T) {
	Convey("Given a DLTask controller", t, func() {
		testInit(t, func(suite *testSuite) {
			Convey("When cancelling with an invalid route (no task ID)", func() {
				var respBody map[string]any

				// Make request without task ID
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodDelete,
					"/dlTask/",
					nil,
					&respBody,
					http.StatusNotFound,
				)

				So(respBody["message"], ShouldEqual, "Endpoint not found")
			})

			Convey("When cancelling with an empty task ID", func() {
				var respBody map[string]any

				// Make request with explicit empty ID parameter
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodDelete,
					"/dlTask/",
					nil,
					&respBody,
					http.StatusNotFound,
				)

				So(respBody["message"], ShouldEqual, "Endpoint not found")
			})

			Convey("When cancelling a non-existent task", func() {
				taskID := "non-existent-task"
				var respBody map[string]interface{}

				// Mock TaskManager to return error for non-existent task
				suite.taskManager.EXPECT().
					GetTaskProgress(taskID).
					Return(int64(0), errors.New("task not found")).
					Times(1)

				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodDelete,
					"/dlTask/"+taskID,
					nil,
					&respBody,
					http.StatusNotFound,
				)

				So(respBody["error"], ShouldEqual, "task not found")
			})

			Convey("When cancelling an already completed task", func() {
				taskID := "completed-task"
				var respBody map[string]interface{}

				// Mock TaskManager to return 100% progress (completed)
				suite.taskManager.EXPECT().
					GetTaskProgress(taskID).
					Return(int64(100), nil).
					Times(1)

				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodDelete,
					"/dlTask/"+taskID,
					nil,
					&respBody,
					http.StatusBadRequest,
				)

				So(respBody["error"], ShouldEqual, "task already completed")
			})

			Convey("When cancelling an already cancelled task", func() {
				taskID := "already-cancelled-task"
				var respBody map[string]interface{}

				// Mock TaskManager to return -1 progress (cancelled)
				suite.taskManager.EXPECT().
					GetTaskProgress(taskID).
					Return(int64(-1), nil).
					Times(1)

				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodDelete,
					"/dlTask/"+taskID,
					nil,
					&respBody,
					http.StatusBadRequest,
				)

				So(respBody["error"], ShouldEqual, "task already cancelled")
			})

			Convey("When successfully cancelling a task", func() {
				taskID := "active-task"
				currentProgress := int64(50)
				var respBody map[string]interface{}

				// Mock TaskManager to return active progress
				suite.taskManager.EXPECT().
					GetTaskProgress(taskID).
					Return(currentProgress, nil).
					Times(1)

				// Mock TaskManager to successfully cancel the task
				suite.taskManager.EXPECT().
					CancelTask(taskID).
					Return(nil).
					Times(1)

				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodDelete,
					"/dlTask/"+taskID,
					nil,
					&respBody,
					http.StatusOK,
				)

				So(respBody["task_id"], ShouldEqual, taskID)
				So(respBody["status_before_cancel"], ShouldEqual, float64(currentProgress)) // JSON numbers become float64
				So(respBody["status"], ShouldEqual, "task cancelled")
			})

			Convey("When cancellation fails due to an error", func() {
				taskID := "error-task"
				currentProgress := int64(30)
				var respBody map[string]interface{}

				// Mock TaskManager to return active progress
				suite.taskManager.EXPECT().
					GetTaskProgress(taskID).
					Return(currentProgress, nil).
					Times(1)

				// Mock TaskManager to return an error during cancellation
				suite.taskManager.EXPECT().
					CancelTask(taskID).
					Return(errors.New("cancellation error")).
					Times(1)

				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodDelete,
					"/dlTask/"+taskID,
					nil,
					&respBody,
					http.StatusNotFound,
				)

				So(respBody["error"], ShouldEqual, "task not found")
			})
		})
	})
}
