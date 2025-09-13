package filters

import (
	"image"
	"image/color"
)

// image to grayscale
func Grayscale(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewGray(bounds)

	// processing each pixel
	for y := bounds.Min.Y; y <bounds.Max.Y; y++{
		for x  := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := src.At(x, y)
			r, g,b, _ := originalColor.RGBA()

			// convertin to 8bit values
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			gray := uint8(0.299*float64(r8) + 0.587*float64(g8) +0.114*float64(b8))
			
			// setting grayscale pixel
			dst.Set(x, y, color.Gray{Y: gray})
		}
	}

	return dst
}