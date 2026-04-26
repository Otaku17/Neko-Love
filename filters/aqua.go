package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Aqua applies an "aqua" themed filter to the provided image.
// The filter analyzes the luminance of each pixel and maps it to a specific color palette
// that emphasizes blue and cyan tones, creating an aquatic effect.
// The resulting image preserves the original alpha channel.
//
// Parameters:
//
//	img image.Image - The source image to apply the filter to.
//
// Returns:
//
//	image.Image - A new image with the aqua filter applied.
func Aqua(img image.Image) image.Image {
	return applyPaletteLuminanceFilter(img, func(lum uint8) (uint8, uint8, uint8) {
		switch {
		case lum >= 235:
			return 255, 255, 255
		case lum >= 179:
			return 80, 220, 255
		case lum >= 115:
			return 15, 100, 120
		case lum >= 38:
			return 35, 39, 42
		default:
			return 35, 39, 42
		}
	})
}
