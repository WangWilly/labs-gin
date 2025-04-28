package taskmanager

import (
	"context"
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////

type Config struct {
	NumWorkers int `env:"NUM_WORKERS,default=4"`
}

type TaskPool struct {
	maxWorkers int
	tasks      chan Task
	idTaskMap  map[string]Task
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

////////////////////////////////////////////////////////////////////////////////

func NewTaskPool(cfg Config) *TaskPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskPool{
		maxWorkers: cfg.NumWorkers,
		tasks:      make(chan Task, cfg.NumWorkers*10),
		idTaskMap:  make(map[string]Task),
		wg:         sync.WaitGroup{},
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (p *TaskPool) GetCtx() context.Context {
	return p.ctx
}

func (p *TaskPool) SubmitTask(task Task) {
	if task == nil {
		return
	}
	if _, ok := p.idTaskMap[task.GetID()]; ok {
		return
	}
	p.idTaskMap[task.GetID()] = task
	p.tasks <- task
}

func (p *TaskPool) GetTaskProgress(taskID string) (int64, error) {
	if task, ok := p.idTaskMap[taskID]; ok {
		return task.GetProgress(), nil
	}
	return 0, fmt.Errorf("task not found")
}

func (p *TaskPool) CancelTask(taskID string) error {
	if task, ok := p.idTaskMap[taskID]; ok {
		task.Cancel()
		delete(p.idTaskMap, taskID)
		return nil
	}
	return fmt.Errorf("task not found")
}

func (p *TaskPool) Run() {
	for range p.maxWorkers {
		go p.createWorker()
	}
}

func (p *TaskPool) createWorker() {
	p.wg.Add(1)
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				fmt.Println("Task channel closed, exiting worker")
				return
			}
			task.Execute()
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func (p *TaskPool) ShutdownNow() {
	p.cancelFunc()
	close(p.tasks)
	p.wg.Wait()
}
