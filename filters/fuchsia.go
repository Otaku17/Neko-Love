package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Fuchsia applies a custom "fuchsia" color filter to the given image.
// The filter maps the luminance of each pixel to a specific color palette
// with fuchsia and dark tones, creating a stylized effect. The alpha channel
// is preserved from the original image.
//
// Parameters:
//
//	img image.Image - The source image to be filtered.
//
// Returns:
//
//	image.Image - A new image with the fuchsia filter applied.
func Fuchsia(img image.Image) image.Image {
	return applyPaletteLuminanceFilter(img, func(lum uint8) (uint8, uint8, uint8) {
		switch {
		case lum > 235:
			return 255, 255, 255
		case lum > 179:
			return 192, 88, 168
		case lum > 115:
			return 152, 40, 128
		case lum >= 38:
			return 35, 39, 42
		default:
			return 35, 39, 42
		}
	})
}
