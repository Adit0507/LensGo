package models

import "image"

// job to processed by workers
type ProcessingJob struct {
	ID         string
	Image      image.Image
	Operations []Operation
	InputPath  string
	OutputPath string
	ResultChan chan ProcessingResult
}

type Operation struct {
	Type   OperationType
	Params map[string]interface{}
}

type OperationType string

// result of processing job
type ProcessingResult struct {
	Success    bool
	Error      error
	OutputPath string
}

const (
	OpResize    OperationType = "resize"
	OpGrayScale OperationType = "grayscale"
	OpBlur      OperationType = "blur"
)

// requets payload for processing
type ProcessRequest struct {
	Filename   string `json:"filename"`
	Operations []struct {
		Type   string                 `json:"type"`
		Params map[string]interface{} `json:"params,omitempty"`
	} `json:"operations"`
}

// response after processing
type ProcessResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	DownloadURL string `json:"download_url,omitempty"`
}
