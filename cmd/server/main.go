package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Adit0507/image-processing-tool/internal/config"
	"github.com/Adit0507/image-processing-tool/internal/handlers"
	"github.com/Adit0507/image-processing-tool/internal/services"
)

func main() {
	cfg := config.New()

	uploadDir := filepath.Join("web", "static", "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	workerPool := services.NewWorkerPool(runtime.NumCPU())
	workerPool.Start()
	defer workerPool.Stop()

	h := handlers.New(cfg, workerPool)

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/upload", h.Upload)
	http.HandleFunc("/process", h.Process)
	http.HandleFunc("/download/", h.Download)

	// servin static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Worker pool initialized with %d workers", runtime.NumCPU())
	
	if err := http.ListenAndServe(":" + cfg.Port, nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}