package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Deepfry applies a "deep-fry" effect to the given image by increasing saturation,
// contrast, and shifting the color balance towards red-orange tones. This effect
// is achieved by manipulating the RGB channels of each pixel, resulting in a
// visually exaggerated and stylized image. The function returns a new image with
// the applied effect, preserving the original image's dimensions.
func Deepfry(img image.Image) image.Image {
	src := ensureRGBA(img)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		srcOffset := src.PixOffset(bounds.Min.X, y)
		dstOffset := dst.PixOffset(bounds.Min.X, y)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r := int(src.Pix[srcOffset])
			g := int(src.Pix[srcOffset+1])
			b := int(src.Pix[srcOffset+2])

			dst.Pix[dstOffset] = clamp8((r*18)/10 + 50)
			dst.Pix[dstOffset+1] = clamp8((g * 14) / 10)
			dst.Pix[dstOffset+2] = clamp8((b * 8) / 10)
			dst.Pix[dstOffset+3] = src.Pix[srcOffset+3]
			srcOffset += 4
			dstOffset += 4
		}
	}

	return dst
}
