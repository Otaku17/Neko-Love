package filters

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Vaporwave applies a "vaporwave" color filter effect to the given image.
// This effect boosts the red and blue (rose and cyan) channels, slightly reduces the green channel,
// and adds a pinkish-cyan tint to the image, reminiscent of the vaporwave aesthetic.
// The function returns a new image.Image with the filter applied.
//
// Parameters:
//
//	img image.Image - The source image to apply the vaporwave filter to.
//
// Returns:
//
//	image.Image - A new image with the vaporwave filter effect applied.
func Vaporwave(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			newR := clamp8(int(r8*1.2 + 30))
			newG := clamp8(int(g8 * 0.9))
			newB := clamp8(int(b8*1.2 + 20))

			dst.Set(x, y, color.NRGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
		}
	}

	return dst
}
