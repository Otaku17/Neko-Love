package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"

	_ "github.com/chai2010/webp"
)

// Glitch applies a "glitch" visual effect to the given image by randomly shifting the red, green, and blue channels
// horizontally on each scanline, and by adding several random horizontal color bands. This creates a distorted,
// glitch-art appearance. The function returns a new image with the effect applied.
//
// Parameters:
//
//	img image.Image - The source image to which the glitch effect will be applied.
//
// Returns:
//
//	image.Image - A new image with the glitch effect applied.
func Glitch(img image.Image) image.Image {
	src := ensureRGBA(img)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	height := bounds.Dy()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		offsetR := rand.Intn(6) - 3
		offsetG := rand.Intn(6) - 3
		offsetB := rand.Intn(6) - 3

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rx := clamp(x+offsetR, bounds.Min.X, bounds.Max.X-1)
			gx := clamp(x+offsetG, bounds.Min.X, bounds.Max.X-1)
			bx := clamp(x+offsetB, bounds.Min.X, bounds.Max.X-1)

			dstOffset := dst.PixOffset(x, y)
			rOffset := src.PixOffset(rx, y)
			gOffset := src.PixOffset(gx, y)
			bOffset := src.PixOffset(bx, y)

			dst.Pix[dstOffset] = src.Pix[rOffset]
			dst.Pix[dstOffset+1] = src.Pix[gOffset+1]
			dst.Pix[dstOffset+2] = src.Pix[bOffset+2]
			dst.Pix[dstOffset+3] = src.Pix[bOffset+3]
		}
	}

	for i := 0; i < 5; i++ {
		yStart := rand.Intn(height)
		bandHeight := rand.Intn(10) + 5
		colorShift := uint8(rand.Intn(100))

		for y := yStart; y < yStart+bandHeight && y < bounds.Max.Y; y++ {
			rowOffset := dst.PixOffset(bounds.Min.X, y)
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				dst.Pix[rowOffset] ^= colorShift
				dst.Pix[rowOffset+1] ^= colorShift
				dst.Pix[rowOffset+2] ^= colorShift
				rowOffset += 4
			}
		}
	}

	return dst
}
