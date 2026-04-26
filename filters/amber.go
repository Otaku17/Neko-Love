package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Amber applies an "amber" color filter to the provided image.Image and returns a new image.Image.
// The filter maps each pixel's luminance to a specific color palette to create an amber-toned effect.
// The output image preserves the alpha channel of the original image.
//
// The mapping is as follows based on luminance:
//   - lum > 0.92: white (255, 255, 255)
//   - lum > 0.7: amber (255, 191, 73)
//   - lum > 0.45: brown (120, 70, 30)
//   - lum >= 0.15: dark gray (35, 39, 42)
//   - otherwise: dark gray (35, 39, 42)
//
// Parameters:
//
//	img image.Image: The source image to apply the filter to.
//
// Returns:
//
//	image.Image: A new image with the amber filter applied.
func Amber(img image.Image) image.Image {
	return applyPaletteLuminanceFilter(img, func(lum uint8) (uint8, uint8, uint8) {
		switch {
		case lum > 235:
			return 255, 255, 255
		case lum > 179:
			return 255, 191, 73
		case lum > 115:
			return 120, 70, 30
		case lum >= 38:
			return 35, 39, 42
		default:
			return 35, 39, 42
		}
	})
}
