package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Greyscale converts the given image to greyscale using standard luminance calculation.
// It iterates over each pixel, computes the luminance based on the RGB values, and sets
// the resulting pixel to a shade of grey with the original alpha value preserved.
//
// Parameters:
//
//	img image.Image - The source image to be converted to greyscale.
//
// Returns:
//
//	image.Image - A new image in greyscale.
func Greyscale(img image.Image) image.Image {
	src := ensureRGBA(img)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		srcOffset := src.PixOffset(bounds.Min.X, y)
		dstOffset := dst.PixOffset(bounds.Min.X, y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			lum := luminance8(src.Pix[srcOffset], src.Pix[srcOffset+1], src.Pix[srcOffset+2])

			dst.Pix[dstOffset] = lum
			dst.Pix[dstOffset+1] = lum
			dst.Pix[dstOffset+2] = lum
			dst.Pix[dstOffset+3] = src.Pix[srcOffset+3]
			srcOffset += 4
			dstOffset += 4
		}
	}

	return dst
}
