package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Blurple applies a "blurple" (blue-purple) themed filter to the given image.
// The filter maps each pixel's luminance to a specific color in a blurple palette,
// producing a stylized effect reminiscent of certain branding themes (e.g., Discord).
// The output image preserves the original alpha channel.
//
// Parameters:
//
//	img image.Image - The source image to be filtered.
//
// Returns:
//
//	image.Image - A new image with the blurple filter applied.
func Blurple(img image.Image) image.Image {
	return applyPaletteLuminanceFilter(img, func(lum uint8) (uint8, uint8, uint8) {
		switch {
		case lum >= 235:
			return 255, 255, 255
		case lum >= 179:
			return 88, 101, 242
		case lum >= 115:
			return 69, 79, 191
		case lum >= 38:
			return 35, 39, 42
		default:
			return 35, 39, 42
		}
	})
}
