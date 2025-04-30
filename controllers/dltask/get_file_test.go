package dltask

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

func TestGetFile(t *testing.T) {
	Convey("Given a DLTask controller", t, func() {
		testInit(t, func(suite *testSuite) {
			// Create a temporary test directory
			testDir := filepath.Join(os.TempDir(), "get_file_test")
			err := os.MkdirAll(testDir, 0755)
			So(err, ShouldBeNil)
			defer os.RemoveAll(testDir)

			// Override the download folder to our test directory
			suite.controller.cfg.DlFolderRoot = testDir

			Convey("When getting a file with an invalid route (no file ID)", func() {
				var respBody map[string]any

				// Make request without file ID
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/dlTaskFile/",
					nil,
					&respBody,
					http.StatusNotFound,
				)

				So(respBody["error"], ShouldEqual, nil)
				So(respBody["message"], ShouldEqual, "Endpoint not found")
			})

			Convey("When getting a file with an invalid file ID (path traversal)", func() {
				var respBody map[string]any

				// Try path traversal
				suite.testServer.MustDoAndMatchCode(
					t,
					http.MethodGet,
					"/dlTaskFile/../../etc/passwd",
					nil,
					&respBody,
					http.StatusNotFound,
				)

				So(respBody["error"], ShouldEqual, nil)
				So(respBody["message"], ShouldEqual, "Endpoint not found")
			})

			Convey("When getting a non-existent file", func() {
				// Request a file that doesn't exist
				resp, err := http.Get(suite.testServer.GetURL(t, "/dlTaskFile/nonexistent.mp4"))
				So(err, ShouldBeNil)
				defer resp.Body.Close()

				So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("When getting an existing file", func() {
				// Create a test file
				fileID := "testfile.mp4"
				filePath := filepath.Join(testDir, fileID)
				testContent := []byte("test video content")
				err := os.WriteFile(filePath, testContent, 0644)
				So(err, ShouldBeNil)

				// Request the file
				resp, err := http.Get(suite.testServer.GetURL(t, "/dlTaskFile/"+fileID))
				So(err, ShouldBeNil)
				defer resp.Body.Close()

				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(resp.Header.Get("Content-Type"), ShouldEqual, "video/mp4")
				So(resp.Header.Get("Accept-Ranges"), ShouldEqual, "bytes")

				// Check content
				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(body, ShouldResemble, testContent)
			})

			Convey("When requesting a range of an existing file", func() {
				// Create a test file with content
				fileID := "rangetest.mp4"
				filePath := filepath.Join(testDir, fileID)
				testContent := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
				err := os.WriteFile(filePath, testContent, 0644)
				So(err, ShouldBeNil)

				// Create a request with Range header
				req, err := http.NewRequest("GET", suite.testServer.GetURL(t, "/dlTaskFile/"+fileID), nil)
				So(err, ShouldBeNil)
				req.Header.Set("Range", "bytes=5-15") // Request bytes 5 through 15

				// Send the request
				client := &http.Client{}
				resp, err := client.Do(req)
				So(err, ShouldBeNil)
				defer resp.Body.Close()

				// Check response
				So(resp.StatusCode, ShouldEqual, http.StatusPartialContent)
				So(resp.Header.Get("Content-Range"), ShouldEqual, "bytes 5-15/26")
				So(resp.Header.Get("Content-Length"), ShouldEqual, "11")

				// Check content
				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(body, ShouldResemble, []byte("FGHIJKLMNOP"))
			})

			Convey("When requesting an invalid range", func() {
				// Create a test file
				fileID := "invalidrange.mp4"
				filePath := filepath.Join(testDir, fileID)
				testContent := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
				err := os.WriteFile(filePath, testContent, 0644)
				So(err, ShouldBeNil)

				// Create a request with invalid Range header
				req, err := http.NewRequest("GET", suite.testServer.GetURL(t, "/dlTaskFile/"+fileID), nil)
				So(err, ShouldBeNil)
				req.Header.Set("Range", "bytes=30-40") // Beyond file size

				// Send the request
				client := &http.Client{}
				resp, err := client.Do(req)
				So(err, ShouldBeNil)
				defer resp.Body.Close()

				// Check response
				So(resp.StatusCode, ShouldEqual, http.StatusRequestedRangeNotSatisfiable)
			})
		})
	})
}
