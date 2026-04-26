package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Posterize applies a posterization effect to the given image by reducing the number of color levels.
// The function processes each pixel, mapping its red, green, and blue channels to one of a fixed number
// of discrete levels (in this case, 4). The result is an image with fewer distinct colors, creating a
// stylized, poster-like appearance.
//
// Parameters:
//
//	img image.Image - The source image to be posterized.
//
// Returns:
//
//	image.Image - A new image with the posterization effect applied.
func Posterize(img image.Image) image.Image {
	src := ensureRGBA(img)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	const levelSize = 64

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		srcOffset := src.PixOffset(bounds.Min.X, y)
		dstOffset := dst.PixOffset(bounds.Min.X, y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Pix[dstOffset] = (src.Pix[srcOffset] / levelSize) * levelSize
			dst.Pix[dstOffset+1] = (src.Pix[srcOffset+1] / levelSize) * levelSize
			dst.Pix[dstOffset+2] = (src.Pix[srcOffset+2] / levelSize) * levelSize
			dst.Pix[dstOffset+3] = src.Pix[srcOffset+3]
			srcOffset += 4
			dstOffset += 4
		}
	}

	return dst
}
