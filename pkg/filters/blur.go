package filters

import (
	"image"
	"image/color"
	"math"
)

func Blur(src image.Image, radius float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	// gaussian kernel
	kernel := createGaussianKernel(radius)
	kernelSize := len(kernel)
	offset := kernelSize / 2

	// applying blur filter
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a float64
			var totalWeight float64

			// convolve with kernel
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					// source pixel coordinates
					sx := x + kx - offset
					sy := y + ky - offset

					// handling boundary conditions by clamping
					if sx < bounds.Min.X {
						sx = bounds.Min.X
					} else if sx >= bounds.Max.X {
						sx = bounds.Max.X - 1
					}
					if sy < bounds.Min.Y {
						sy = bounds.Min.Y
					} else if sy >= bounds.Max.Y {
						sy = bounds.Max.Y - 1
					}

					// pixel color
					srcColor := src.At(sx, sy)
					sr, sg, sb, sa := srcColor.RGBA()

					// Get kernel weight
					weight := kernel[ky][kx]
					totalWeight += weight

					// accumulate weighted color values
					r += float64(sr>>8) * weight
					g += float64(sg>>8) * weight
					b += float64(sb>>8) * weight
					a += float64(sa>>8) * weight
				}
			}

			if totalWeight > 0 {
				r /= totalWeight
				g /= totalWeight
				b /= totalWeight
				a /= totalWeight
			}

			dst.Set(x, y, color.RGBA{
				R: uint8(clamp(r, 0, 255)),
				G: uint8(clamp(g, 0, 255)),
				B: uint8(clamp(b, 0, 255)),
				A: uint8(clamp(a, 0, 255)),
			})
		}
	}

	return dst
}

// generates a gaussian blur
func createGaussianKernel(radius float64) [][]float64 {
	size := int(math.Ceil(radius*2)) + 1
	if size%2 == 0 {
		size++ //ensure odd size
	}

	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
	}

	center := size / 2
	sigma := radius / 3.0
	twoSigmaSquared := 2 * sigma * sigma

	var sum float64
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x - center)
			dy := float64(y - center)

			distance := dx*dx + dy*dy
			value := math.Exp(-distance / twoSigmaSquared)
			kernel[y][x] = value
			sum += value
		}
	}

	// normalize kernel
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= sum
		}
	}

	return kernel
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
