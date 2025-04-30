package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/WangWilly/labs-gin/pkgs/utils"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

// Controller defines the interface that controllers must implement
// to be used with the test HTTP server
type Controller interface {
	RegisterRoutes(router *gin.Engine)
}

////////////////////////////////////////////////////////////////////////////////

type TestHttpServer struct {
	Server *httptest.Server
	client *http.Client
}

func NewTestHttpServer(controller Controller) TestHttpServer {
	router := utils.GetDefaultRouter()
	controller.RegisterRoutes(router)
	server := httptest.NewServer(router)
	client := server.Client()
	return TestHttpServer{
		Server: server,
		client: client,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (c *TestHttpServer) GetURL(t *testing.T, path string) string {
	url, err := url.JoinPath(c.Server.URL, path)
	if err != nil {
		t.Fatal(err)
	}
	return url
}

func (c *TestHttpServer) MustDo(
	t *testing.T,
	method string,
	url string,
	reqBody interface{},
	respBody interface{},
) int {
	// Encode request
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare request
	req, err := http.NewRequest(method, c.GetURL(t, url), &buf)
	if err != nil {
		t.Fatal(err)
	}

	// Do request
	resp, err := c.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse response body
	if respBody != nil {
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			t.Logf("resp.Body: %s", resp.Body)
			t.Fatalf("status: %d, err:%s", resp.StatusCode, err)
		}
	}

	return resp.StatusCode
}

////////////////////////////////////////////////////////////////////////////////

func (c *TestHttpServer) MustDoAndMatchCode(
	t *testing.T,
	method string,
	url string,
	reqBody any,
	respBody any,
	code int,
) {
	respCode := c.MustDo(t, method, url, reqBody, respBody)
	So(respCode, ShouldEqual, code)
}

/**
func (c *TestHttpServer) MustSucceeded(t *testing.T, method string, url string, reqBody interface{}, respBody interface{}, code int) {
	respCode := c.MustDo(t, method, url, reqBody, respBody)
	So(respCode, ShouldEqual, code)
}

func (c *TestHttpServer) MustFailed(t *testing.T, method string, path string, reqBody interface{}, err *errors.Error) {
	respBody := new(errors.Error)
	respCode := c.MustDo(t, method, path, reqBody, respBody)
	So(respBody.Message, ShouldEqual, err.Message)
	So(respCode, ShouldEqual, err.Code)
}
*/
