package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Negative returns a new image that is the negative (color-inverted) version of the input image.
// Each pixel's red, green, and blue channels are inverted, while the alpha channel is preserved.
// The function supports any image.Image input and outputs an *image.RGBA.
func Negative(img image.Image) image.Image {
	src := ensureRGBA(img)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		srcOffset := src.PixOffset(bounds.Min.X, y)
		dstOffset := dst.PixOffset(bounds.Min.X, y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Pix[dstOffset] = 255 - src.Pix[srcOffset]
			dst.Pix[dstOffset+1] = 255 - src.Pix[srcOffset+1]
			dst.Pix[dstOffset+2] = 255 - src.Pix[srcOffset+2]
			dst.Pix[dstOffset+3] = src.Pix[srcOffset+3]
			srcOffset += 4
			dstOffset += 4
		}
	}

	return dst
}
