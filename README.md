# LensGo 🖼️ 

A concurrent image processing tool built in Go, designed to explore parallel processing. It provides a simple web interface for uploading images, applying filters, resizing, and downloading processed results  all powered by a worker pool concurrency model. 


## Demo
https://github.com/user-attachments/assets/e2f8a920-15a0-449e-a6b7-77ef02d0190a


## ✨ Features
- 🖼️ Image Resizing
   - Bilinear
 - 🎨 Filters
   - Blur
   - Grayscale

- 📂 Format Support
  - JPEG
  - PNG
  - WebP
    
- ⚡Performance-Oriented
  - Parallel processing with a worker pool
  - Batch operation support
  - Efficient memory usage for large images
  - Progress tracking for long operations
 
## 🏗️ Architecture Overview
- Worker Pool Design
  - Configurable worker pool (default: runtime.NumCPU())
  - Job queue for processing tasks
  - Results returned via channels
  - Graceful shutdown support

- Processing Flow
  - User uploads image via web form
  - Server validates file (size, format, type)
  - Processing job created → submitted to worker pool
  - Worker processes (resize / filter / chain of ops)
  - Processed image saved temporarily
  - User gets preview + download link

- Web Interface
  - Simple HTML form for uploads
  - AJAX for non-blocking requests
  - Live progress updates
  - Download links for results
 
## Key Design Decisions
- <b>Concurrency Model</b> → Worker pool for controlled CPU usage

- <b>File Handling</b> → Temporary storage with cleanup

- <b>API Design →</b> REST endpoints (/upload, /process, /download)

- <b>Error Handling →</b> Proper error propagation + user feedback

- <b>Security →</b> File validation, size limits, allowed extensions

## 🚀 Getting Started
- Go 1.22+

## Installation
````````````````````
git clone https://github.com/<your-username>/<repo-name>.git
cd <repo-name>
go mod tidy
````````````````````

## Run the server
````````````````````
go run cmd/server/main.go
````````````````````

## Usage

- Open http://localhost:8080 in your browser

- Upload an image

- Choose operations (resize, filters)

- Download results


