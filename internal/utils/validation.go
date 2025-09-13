package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// if uploaded file meets requirement
func ValidateFile(header *multipart.FileHeader, allowedFormats []string, maxSize int64 ) bool{
	if header.Size > maxSize{
		return  false
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	for _, allowed := range allowedFormats{
		if ext == allowed{
			return true
		}
	}

	return true
}

func GenerateUniqueFilename(originalFilename string) string{
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}

func GenerateProcessedFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_processed_%d%s", name, timestamp, ext)
}

func GenerateJobID() string {
	return fmt.Sprintf("job_%d", time.Now().UnixNano())
}