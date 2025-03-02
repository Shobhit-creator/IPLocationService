package workerpool

import (
	"errors"
	"fmt"
)

type Result interface{}

type WorkerPool struct {
	Tasks           chan func()
	NumberOfWorkers int
	quit         chan struct{}
}


func NewWorkerPool(numberOfWorkers int, bufferSize int) *WorkerPool {
	return &WorkerPool{Tasks: make(chan func(), bufferSize), NumberOfWorkers: numberOfWorkers, quit: make(chan struct{})}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.NumberOfWorkers; i++ {
        go wp.worker(i)
    }
}
func (wp *WorkerPool) worker(id int) {
	for {
		select {
		case task, ok := <-wp.Tasks:
			if !ok {
				fmt.Printf("Worker %d stopping: Task channel closed\n", id)
				return
			}
			task()
		case <-wp.quit:
		    fmt.Printf("Worker %d received stop signal\n", id)
            return
		}
	}
}


func (wp *WorkerPool) Submit(task func()) error {
	select {
		case wp.Tasks <- task:
			return nil
		default:
			return errors.New("worker pool is full")
	}
}

func (wp *WorkerPool) Stop() {
    close(wp.quit)
    close(wp.Tasks)
}