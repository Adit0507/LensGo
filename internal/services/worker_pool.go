package services

import (
	"log"
	"sync"

	"github.com/Adit0507/image-processing-tool/internal/models"
)

// manages a pool of workers for concurrent image processing
type WorkerPool struct {
	workerCount int
	jobQueue    chan models.ProcessingJob
	quit        chan bool
	wg          sync.WaitGroup
	processor   *ImageProcessor
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan models.ProcessingJob, workerCount*2),
		quit:        make(chan bool),
		processor:   NewImageProcessor(),
	}
}

func (wp *WorkerPool) Start() {
	log.Printf("Starting worker pools with %d workers", wp.workerCount)

	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i + 1)
	}
}

func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	close(wp.quit)
	wp.wg.Wait()
	close(wp.jobQueue)

	log.Println("Worker pool stopped")
}

func (wp *WorkerPool) worker(id int) {

}
