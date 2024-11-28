package services

import (
	"context"
	"log/slog"
	"sync"
)

const (
	MAX_WORKERS = 10
	TAKS_BUFF   = 10
)

type Task struct {
	ImageID       string
	ImagePath     string
	ImageFormat   string
	ThumbnailPath string
	Width         int
	Height        int
}

type WorkerPool struct {
	tasks chan Task
	wg    sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		tasks: make(chan Task, TAKS_BUFF),
	}
}

func (p *WorkerPool) Start(ctx context.Context, thumbnailFunc func(context.Context, *Task) error) {
	for i := 0; i < MAX_WORKERS; i++ {
		go func() {
			for task := range p.tasks {
				if err := thumbnailFunc(ctx, &task); err != nil {
					slog.Error("Error processing task",
						slog.String("imageID", task.ImageID),
						slog.String("error", err.Error()))
				}
				p.wg.Done()
			}
		}()
	}
}

func (p *WorkerPool) AddTask(task Task) {
	p.wg.Add(1)
	p.tasks <- task
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.tasks)
}
