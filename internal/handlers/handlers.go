package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Adit0507/image-processing-tool/internal/config"
	"github.com/Adit0507/image-processing-tool/internal/models"
	"github.com/Adit0507/image-processing-tool/internal/services"
	"github.com/Adit0507/image-processing-tool/internal/utils"
)

type Handler struct {
	config     *config.Config
	workerPool *services.WorkerPool
}

func New(cfg *config.Config, wp *services.WorkerPool) *Handler {
	return &Handler{
		config:     cfg,
		workerPool: wp,
	}
}

// serving main page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Template execution error %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parsing multipart form
	if err := r.ParseMultipartForm(h.config.MaxFileSize); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// validate file
	if !utils.ValidateFile(header, h.config.AllowedFormats, h.config.MaxFileSize) {
		http.Error(w, "Invalid file format or size", http.StatusBadRequest)
		return
	}

	// creatin unique filename
	filename := utils.GenerateUniqueFilename(header.Filename)
	filePath := filepath.Join(h.config.UploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("error saving file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)

		return
	}

	response := map[string]string{
		"status":   "success",
		"filename": filename,
		"message":  "File uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	inputPath := filepath.Join(h.config.UploadDir, req.Filename)
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// generatte output filenames
	outputFilename := utils.GenerateProcessedFilename(req.Filename)
	outputPath := filepath.Join(h.config.UploadDir, outputFilename)

	// convert operations
	var operations []models.Operation
	for _, op := range req.Operations {
		operations = append(operations, models.Operation{
			Type:   models.OperationType(op.Type),
			Params: op.Params,
		})
	}

	// creatin processing job
	resultChan := make(chan models.ProcessingResult, 1)
	job := models.ProcessingJob{
		ID:         utils.GenerateJobID(),
		InputPath:  inputPath,
		OutputPath: outputPath,
		Operations: operations,
		ResultChan: resultChan,
	}

	h.workerPool.SubmitJob(job)

	// waiitn for result
	select {
	case result := <-resultChan:
		if !result.Success {
			log.Printf("Processing error: %v", result.Error)
			http.Error(w, "Processing failed", http.StatusInternalServerError)
			return
		}

		response := models.ProcessResponse{
			Success:     true,
			Message:     "Image processed successfully",
			DownloadURL: "/download/" + outputFilename,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case <-time.After(30 * time.Second):
		http.Error(w, "Processing timeout", http.StatusRequestTimeout)
	}
}
