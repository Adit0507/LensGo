package utils

import (
	"mime/multipart"
	"path/filepath"
	"strings"
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