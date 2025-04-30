package taskmanager

import (
	"context"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	gomock "go.uber.org/mock/gomock"
)

////////////////////////////////////////////////////////////////////////////////

func TestNewTaskPool(t *testing.T) {
	Convey("Given a configuration", t, func() {
		cfg := Config{
			NumWorkers: 4,
		}

		Convey("When creating a new task pool", func() {
			pool := NewTaskPool(cfg)

			Convey("Then the pool should be properly initialized", func() {
				So(pool, ShouldNotBeNil)
				So(pool.maxWorkers, ShouldEqual, cfg.NumWorkers)
				So(pool.tasks, ShouldNotBeNil)
				So(cap(pool.tasks), ShouldEqual, cfg.NumWorkers*10)
				So(pool.idTaskMap, ShouldNotBeNil)
				So(pool.ctx, ShouldNotBeNil)
				So(pool.cancelFunc, ShouldNotBeNil)
			})
		})
	})
}

func TestGetCtx(t *testing.T) {
	Convey("Given a task pool", t, func() {
		pool := NewTaskPool(Config{NumWorkers: 2})

		Convey("When getting the context", func() {
			ctx := pool.GetCtx()

			Convey("Then it should return the pool's context", func() {
				So(ctx, ShouldEqual, pool.ctx)
			})
		})
	})
}

func TestSubmitTask(t *testing.T) {
	Convey("Given a task pool", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pool := NewTaskPool(Config{NumWorkers: 2})

		Convey("When submitting a nil task", func() {
			initialLength := len(pool.idTaskMap)
			pool.SubmitTask(nil)

			Convey("Then the task should be ignored", func() {
				So(len(pool.idTaskMap), ShouldEqual, initialLength)
			})
		})

		Convey("When submitting a valid task", func() {
			mockTask := NewMockTask(ctrl)
			taskID := "test-task-1"
			mockTask.EXPECT().GetID().Return(taskID).AnyTimes()

			// Submit the task
			pool.SubmitTask(mockTask)

			Convey("Then the task should be added to the pool", func() {
				// Check if task is in the map
				task, exists := pool.idTaskMap[taskID]
				So(exists, ShouldBeTrue)
				So(task, ShouldEqual, mockTask)

				// Check if task is in the channel
				select {
				case taskFromChan := <-pool.tasks:
					So(taskFromChan, ShouldEqual, mockTask)
				default:
					t.Fatal("Task not added to the channel")
				}
			})
		})

		Convey("When submitting a duplicate task", func() {
			mockTask := NewMockTask(ctrl)
			taskID := "test-task-2"
			mockTask.EXPECT().GetID().Return(taskID).AnyTimes()

			// Submit the task first time
			pool.SubmitTask(mockTask)

			// Drain the channel to check it later
			<-pool.tasks

			// Submit the same task again
			pool.SubmitTask(mockTask)

			Convey("Then the duplicate submission should be ignored", func() {
				// Check if channel is empty (no new task added)
				select {
				case <-pool.tasks:
					t.Fatal("Duplicate task was added to the channel")
				default:
					// This is expected - channel should be empty
				}
			})
		})
	})
}

func TestGetTaskProgress(t *testing.T) {
	Convey("Given a task pool with tasks", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pool := NewTaskPool(Config{NumWorkers: 2})

		// Create and add a task
		taskID := "progress-test-task"
		expectedProgress := int64(50)

		mockTask := NewMockTask(ctrl)
		mockTask.EXPECT().GetID().Return(taskID).AnyTimes()
		mockTask.EXPECT().GetProgress().Return(expectedProgress).AnyTimes()

		pool.idTaskMap[taskID] = mockTask

		Convey("When getting progress for an existing task", func() {
			progress, err := pool.GetTaskProgress(taskID)

			Convey("Then it should return the correct progress", func() {
				So(err, ShouldBeNil)
				So(progress, ShouldEqual, expectedProgress)
			})
		})

		Convey("When getting progress for a non-existent task", func() {
			progress, err := pool.GetTaskProgress("non-existent-task")

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "task not found")
				So(progress, ShouldEqual, 0)
			})
		})
	})
}

func TestCancelTask(t *testing.T) {
	Convey("Given a task pool with tasks", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pool := NewTaskPool(Config{NumWorkers: 2})

		// Create and add a task
		taskID := "cancel-test-task"

		mockTask := NewMockTask(ctrl)

		pool.idTaskMap[taskID] = mockTask

		Convey("When cancelling an existing task", func() {
			mockTask.EXPECT().GetID().Return(taskID).AnyTimes()
			mockTask.EXPECT().Cancel().Times(1)

			err := pool.CancelTask(taskID)

			Convey("Then it should cancel the task and remove it from the map", func() {
				So(err, ShouldBeNil)
				_, exists := pool.idTaskMap[taskID]
				So(exists, ShouldBeFalse)
			})
		})

		Convey("When cancelling a non-existent task", func() {
			err := pool.CancelTask("non-existent-task")

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "task not found")
			})
		})
	})
}

func TestRunAndWorkers(t *testing.T) {
	Convey("Given a task pool", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pool := NewTaskPool(Config{NumWorkers: 2})

		Convey("When running the worker pool with a successful task", func() {
			mockTask := NewMockTask(ctrl)
			taskID := "successful-task"

			mockTask.EXPECT().GetID().Return(taskID).AnyTimes()
			mockTask.EXPECT().Execute().Return(true).Times(1)

			// Start the pool
			pool.Run()

			// Submit the task
			pool.SubmitTask(mockTask)

			// Give some time for workers to process
			time.Sleep(100 * time.Millisecond)

			Convey("Then the task should be executed successfully", func() {
				// Successful execution doesn't remove the task from the map
				// in the current implementation
				_, exists := pool.idTaskMap[taskID]
				So(exists, ShouldBeTrue)
			})

			// Clean up
			pool.ShutdownNow()
		})

		Convey("When running the worker pool with a failing task that needs retry", func() {
			mockTask := NewMockTask(ctrl)
			taskID := "failing-task"
			retryChan := make(chan struct{}, 1)

			mockTask.EXPECT().GetID().Return(taskID).AnyTimes()
			mockTask.EXPECT().Execute().Return(false).Times(1)
			mockTask.EXPECT().SetRetrySignal().Return(retryChan).Times(1)

			// Start the pool
			pool.Run()

			// Submit the task
			pool.SubmitTask(mockTask)

			// Give some time for workers to process
			time.Sleep(100 * time.Millisecond)

			Convey("Then the task should be removed from the map and retry signal requested", func() {
				_, exists := pool.idTaskMap[taskID]
				So(exists, ShouldBeFalse)
			})

			// Clean up
			pool.ShutdownNow()
		})
	})
}

func TestShutdownNow(t *testing.T) {
	Convey("Given a running task pool", t, func() {
		pool := NewTaskPool(Config{NumWorkers: 2})
		pool.Run()

		Convey("When shutting down the pool", func() {
			// Use a WaitGroup to ensure we can detect the end of ShutdownNow
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				pool.ShutdownNow()
				wg.Done()
			}()

			// Use a timeout for the wait
			waitCh := make(chan struct{})
			go func() {
				wg.Wait()
				close(waitCh)
			}()

			Convey("Then it should terminate gracefully", func() {
				select {
				case <-waitCh:
					// ShutdownNow completed successfully
					So(true, ShouldBeTrue)
				case <-time.After(1 * time.Second):
					t.Fatal("Timeout waiting for ShutdownNow to complete")
				}

				// Check that the context is canceled
				select {
				case <-pool.ctx.Done():
					So(true, ShouldBeTrue)
				default:
					t.Fatal("Context was not canceled")
				}
			})
		})
	})
}

func TestTaskRetryAndResubmission(t *testing.T) {
	Convey("Given a task pool", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pool := NewTaskPool(Config{NumWorkers: 2})
		pool.Run()
		defer pool.ShutdownNow()

		mockTask := NewMockTask(ctrl)
		taskID := "retry-task"
		retryChan := make(chan struct{}, 1)

		mockTask.EXPECT().GetID().Return(taskID).MinTimes(1)
		mockTask.EXPECT().Execute().Return(false).Times(1)
		mockTask.EXPECT().SetRetrySignal().Return(retryChan).Times(1)

		// For resubmission verification
		var resubmissionCheck sync.WaitGroup
		resubmissionCheck.Add(1)

		secondExecuteCall := mockTask.EXPECT().GetID().Return(taskID).MaxTimes(1)
		mockTask.EXPECT().Execute().Return(true).Times(1).After(secondExecuteCall).Do(func() {
			resubmissionCheck.Done()
		})

		// Submit the task
		pool.SubmitTask(mockTask)

		// Give some time for initial execution
		time.Sleep(50 * time.Millisecond)

		Convey("When the task fails and needs to be retried", func() {
			// Trigger the retry signal
			retryChan <- struct{}{}

			// Wait for resubmission with timeout
			waitCh := make(chan struct{})
			go func() {
				resubmissionCheck.Wait()
				close(waitCh)
			}()

			Convey("Then the task should be resubmitted and executed again", func() {
				select {
				case <-waitCh:
					// Resubmission completed successfully
					So(true, ShouldBeTrue)
				case <-time.After(1 * time.Second):
					t.Fatal("Timeout waiting for task resubmission and execution")
				}
			})
		})
	})
}

func TestTaskPoolWithContextCancellation(t *testing.T) {
	Convey("Given a task pool with tasks", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create a parent context that we can cancel
		parentCtx, cancelParent := context.WithCancel(context.Background())

		// Create a custom constructor to use our parent context
		customNewTaskPool := func(cfg Config) *TaskPool {
			ctx, cancel := context.WithCancel(parentCtx)
			return &TaskPool{
				maxWorkers: cfg.NumWorkers,
				tasks:      make(chan Task, cfg.NumWorkers*10),
				idTaskMap:  make(map[string]Task),
				wg:         sync.WaitGroup{},
				ctx:        ctx,
				cancelFunc: cancel,
			}
		}

		pool := customNewTaskPool(Config{NumWorkers: 2})

		Convey("When the parent context is cancelled", func() {
			// Add a WaitGroup to detect when the workers exit
			var workerExitWg sync.WaitGroup

			// Replace createWorker to notify when it exits due to context cancellation
			createWorker := func() {
				workerExitWg.Add(1)
				pool.wg.Add(1)
				defer pool.wg.Done()

				for {
					select {
					case <-pool.ctx.Done():
						workerExitWg.Done()
						return
					case task, ok := <-pool.tasks:
						if !ok {
							return
						}
						_ = task.Execute()
					}
				}
			}
			pool.runFor(createWorker)

			// Cancel the parent context
			cancelParent()

			// Wait for the worker to detect the cancellation
			waitCh := make(chan struct{})
			go func() {
				workerExitWg.Wait()
				close(waitCh)
			}()

			Convey("Then the workers should exit", func() {
				select {
				case <-waitCh:
					// Worker exited due to context cancellation
					So(true, ShouldBeTrue)
				case <-time.After(1 * time.Second):
					t.Fatal("Timeout waiting for worker to exit after context cancellation")
				}
			})
		})
	})
}
