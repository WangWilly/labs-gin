package dltask

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

func TestGetStatus(t *testing.T) {
	Convey("Given a DLTask controller", t, func() {
		testInit(t, func(suite *testSuite) {
			// TODO:
			Convey("When getting status with an invalid route (no task ID)", func() {
				// Arrange
				var respBody map[string]any

				// Act
				// We need to test the 404 directly because the router won't match without a parameter
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/dlTask/", // The trailing slash will be the empty tid parameter
					nil,
					&respBody,
					http.StatusNotFound,
				)

				// Assert
				So(respBody["message"], ShouldEqual, "Endpoint not found")
			})

			// TODO:
			Convey("When getting status with an empty task ID", func() {
				// Arrange
				var respBody map[string]any

				// Act - use an explicit empty string as the tid parameter
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/dlTask/", // The trailing slash will be the empty tid parameter
					nil,
					&respBody,
					http.StatusNotFound,
				)

				// Assert
				So(respBody["message"], ShouldEqual, "Endpoint not found")
			})

			Convey("When getting status for a non-existent task", func() {
				// Arrange
				taskID := "non-existent-task"
				var respBody map[string]interface{}

				// Mock TaskManager to return error for non-existent task
				suite.taskManager.EXPECT().
					GetTaskProgress(taskID).
					Return(int64(0), errors.New("task not found")).
					Times(1)

				// Act
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/dlTask/"+taskID,
					nil,
					&respBody,
					http.StatusNotFound,
				)

				// Assert
				So(respBody["error"], ShouldEqual, "task not found")
			})

			Convey("When getting status for an existing task", func() {
				// Arrange
				taskID := "existing-task"
				expectedProgress := int64(75)
				var respBody map[string]interface{} // Changed to map to match JSON structure

				// Mock TaskManager to return progress for existing task
				suite.taskManager.EXPECT().
					GetTaskProgress(taskID).
					Return(expectedProgress, nil).
					Times(1)

				// Act
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/dlTask/"+taskID,
					nil,
					&respBody,
					http.StatusOK,
				)

				// Assert
				So(respBody["task_id"], ShouldEqual, taskID)
				So(respBody["status"], ShouldEqual, float64(expectedProgress)) // JSON numbers become float64
			})
		})
	})
}
