package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID int
}

type WorkerPool struct {
	tasks chan Task
	wg    sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
	wp := &WorkerPool{
		tasks: make(chan Task, size),
	}

	wp.wg.Add(size)
	for i := range size {
		go func(i int) {
			defer wp.wg.Done()

			for t := range wp.tasks {
				fmt.Printf("Task is processing, %v with worker %d\n", t, i)
				time.Sleep(100 * time.Microsecond)
			}
		}(i)
	}
	return wp
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Wait() {
	close(wp.tasks)
	wp.wg.Wait()
}

func main() {

	wp := NewWorkerPool(100)
	for i := range 10000 {
		fmt.Printf("Add %d task\n", i)
		wp.AddTask(Task{ID: i})
	}

	wp.Wait()

}
