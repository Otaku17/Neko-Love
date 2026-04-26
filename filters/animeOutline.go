package filters

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// AnimeOutline applies a simple edge detection filter to the given image.
// It highlights the outlines by comparing the color differences between each pixel
// and its right and bottom neighbors. If the combined color difference exceeds a
// specified threshold, the pixel is set to black (indicating an edge); otherwise,
// the original color is retained. The result is a new image with emphasized outlines,
// suitable for creating an "anime-style" outline effect.
//
// Parameters:
//
//	img image.Image - The source image to process.
//
// Returns:
//
//	image.Image - A new image with detected outlines.
func AnimeOutline(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	threshold := 30 // sensibilité des contours

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			c1 := img.At(x, y)
			c2 := img.At(x+1, y)
			c3 := img.At(x, y+1)

			r1, g1, b1, _ := c1.RGBA()
			r2, g2, b2, _ := c2.RGBA()
			r3, g3, b3, _ := c3.RGBA()

			delta1 := absDiff(r1, r2) + absDiff(g1, g2) + absDiff(b1, b2)
			delta2 := absDiff(r1, r3) + absDiff(g1, g3) + absDiff(b1, b3)

			if delta1+delta2 > uint32(threshold*3*256) {
				dst.Set(x, y, color.Black)
			} else {
				dst.Set(x, y, c1)
			}
		}
	}

	return dst
}

// absDiff returns the absolute difference between two uint32 values a and b.
func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}
