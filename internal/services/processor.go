package services

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"strings"

	"github.com/Adit0507/image-processing-tool/pkg/filters"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// loading image from specified path
func (p *ImageProcessor) LoadImage(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (p *ImageProcessor) SaveImage(img image.Image, filePath string) error{
	file, err := os.Create(filePath)
	if err !=nil{
		return err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext{
		case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case ".png":
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	}
} 

func (p *ImageProcessor) Resize(img image.Image, width, height int) image.Image {
	return filters.Resize(img, width, height)
}

func (p *ImageProcessor) Grayscale(img image.Image) image.Image {
	return filters.Grayscale(img)
}

func (p *ImageProcessor) Blur(img image.Image, radius float64) image.Image {
	return filters.Blur(img, radius)
}