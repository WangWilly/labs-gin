package dltask

import (
	"testing"

	"github.com/WangWilly/labs-gin/pkgs/testutils"
	"github.com/WangWilly/labs-gin/pkgs/uuid"
	gomock "go.uber.org/mock/gomock"
)

////////////////////////////////////////////////////////////////////////////////

type testSuite struct {
	taskManager *MockTaskManager
	uuidGen     *uuid.MockUUID

	controller *Controller
	testServer testutils.TestHttpServer
}

func testInit(t *testing.T, test func(*testSuite)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskManager := NewMockTaskManager(ctrl)
	uuidGen := uuid.NewMockUUID(ctrl)

	cfg := Config{
		DlFolderRoot: "./public/downloads",
		RetryDelay:   5,
		MaxRetries:   3,
		MaxTimeout:   5,
	}
	controller := NewController(cfg, taskManager, uuidGen)
	testServer := testutils.NewTestHttpServer(controller)

	suite := &testSuite{
		taskManager: taskManager,
		uuidGen:     uuidGen,
		controller:  controller,
		testServer:  testServer,
	}

	test(suite)
}
