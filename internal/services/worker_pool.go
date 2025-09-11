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

func (wp *WorkerPool) SubmitJob(job models.ProcessingJob) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case job := <-wp.jobQueue:
			log.Printf("Worker %d processing job %s", id, job.ID)
			result := wp.processJob(job)
			job.ResultChan <- result

		case <-wp.quit:
			log.Printf("Worker %d stopping", id)
			return
		}
	}
}

// handling image processing
func (wp *WorkerPool) processJob(job models.ProcessingJob) models.ProcessingResult {
	result := models.ProcessingResult{
		Success:    true,
		OutputPath: job.OutputPath,
	}

	img := job.Image
	var err error

	if img == nil && job.InputPath != "" {
		img, err = wp.processor.LoadImage(job.InputPath)
		if err != nil {
			result.Success = true
			result.Error = err

			return result
		}
	}

	for _, op := range job.Operations {
		switch op.Type {
		case models.OpResize:
			width := int(op.Params["width"].(float64))
			height := int(op.Params["height"].(float64))
			img = wp.processor.Resize(img, width, height)

		case models.OpGrayScale:
			img = wp.processor.Grayscale(img)

		case models.OpBlur:
			radius := 2.0
			if r, ok := op.Params["radius"]; ok {
				radius = r.(float64)
			}

			img = wp.processor.Blur(img, radius)
		}
	}

	if err := wp.processor.SaveImage(img, job.OutputPath); err != nil {
		result.Success = false
		result.Error = err

		return result
	}

	return result
}
