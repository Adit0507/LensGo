package filters

import (
	"image"
	"image/color"
)

func Resize(src image.Image, newWidth, newHeight int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Max.X - srcBounds.Min.X
	srcHeight := srcBounds.Max.Y - srcBounds.Min.Y

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// scaling ratioos
	xRatio := float64(srcWidth) / float64(newWidth)
	yRatio := float64(srcHeight) / float64(newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := float64(x) * xRatio
			srcY := float64(y) * yRatio

			// 4 nearest pixels for bilinear interpolation
			x1 := int(srcX)
			y1 := int(srcY)
			x2 := x1 + 1
			y2 := y1 + 1

			// Clamp coordinates to image bounds
			x1 = clampInt(x1, srcBounds.Min.X, srcBounds.Max.X-1)
			y1 = clampInt(y1, srcBounds.Min.Y, srcBounds.Max.Y-1)
			x2 = clampInt(x2, srcBounds.Min.X, srcBounds.Max.X-1)
			y2 = clampInt(y2, srcBounds.Min.Y, srcBounds.Max.Y-1)

			// colors of four pixels
			c1 := src.At(x1, y1)
			c2 := src.At(x2, y1)
			c3 := src.At(x1, y2)
			c4 := src.At(x2, y2)

			// Calculate interpolation weights
			wx := srcX - float64(x1)
			wy := srcY - float64(y1)

			// Perform bilinear interpolation
			interpolatedColor := bilinearInterpolate(c1, c2, c3, c4, wx, wy)
			dst.Set(x, y, interpolatedColor)
		}
	}

	return dst
}

func bilinearInterpolate(c1, c2, c3, c4 color.Color, wx, wy float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r3, g3, b3, a3 := c3.RGBA()
	r4, g4, b4, a4 := c4.RGBA()

	// convert to 8bit values
	r1, g1, b1, a1 = r1>>8, g1>>8, b1>>8, a1>>8
	r2, g2, b2, a2 = r2>>8, g2>>8, b2>>8, a2>>8
	r3, g3, b3, a3 = r3>>8, g3>>8, b3>>8, a3>>8
	r4, g4, b4, a4 = r4>>8, g4>>8, b4>>8, a4>>8

	// interpolate horizontally
	r12 := lerp(float64(r1), float64(r2), wx)
	g12 := lerp(float64(g1), float64(g2), wx)
	b12 := lerp(float64(b1), float64(b2), wx)
	a12 := lerp(float64(a1), float64(a2), wx)

	r34 := lerp(float64(r3), float64(r4), wx)
	g34 := lerp(float64(g3), float64(g4), wx)
	b34 := lerp(float64(b3), float64(b4), wx)
	a34 := lerp(float64(a3), float64(a4), wx)

	// interpolate vertically
	r := lerp(r12, r34, wy)
	g := lerp(g12, g34, wy)
	b := lerp(b12, b34, wy)
	a := lerp(a12, a34, wy)

	return color.RGBA{
		R: uint8(clamp(r, 0, 255)),
		G: uint8(clamp(g, 0, 255)),
		B: uint8(clamp(b, 0, 255)),
		A: uint8(clamp(a, 0, 255)),
	}
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
