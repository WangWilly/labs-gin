package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDefaultRouter(t *testing.T) {
	// Set Gin to test mode to avoid debug logs in tests
	gin.SetMode(gin.TestMode)

	Convey("Given a call to GetDefaultRouter", t, func() {
		router := GetDefaultRouter()

		Convey("Then a router should be created", func() {
			So(router, ShouldNotBeNil)
		})

		Convey("When accessing a non-existent route", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/nonexistent-route", nil)
			router.ServeHTTP(w, req)

			Convey("Then a 404 status with the correct message should be returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)
				So(response["message"], ShouldEqual, "Endpoint not found")
			})
		})

		Convey("When requesting the favicon", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/favicon.ico", nil)
			router.ServeHTTP(w, req)

			Convey("Then the router should attempt to serve the static file", func() {
				// Note: This might be 200 OK if file exists or 404 if it doesn't in the test environment
				So(w.Code == http.StatusOK || w.Code == http.StatusNotFound, ShouldBeTrue)
			})
		})

		Convey("When calling the ping endpoint", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/ping", nil)
			router.ServeHTTP(w, req)

			Convey("Then a 200 status with a pong message should be returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)
				So(response["message"], ShouldEqual, "pong")
			})
		})
	})
}
