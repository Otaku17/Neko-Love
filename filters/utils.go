package filters

import (
	"image"
	"image/color"
	"image/draw"
)

// clamp restricts the integer value v to be within the range [min, max].
// If v is less than min, min is returned. If v is greater than max, max is returned.
// Otherwise, v is returned unchanged.
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// clamp8 limits the input integer v to the range [0, 255] and returns it as a uint8.
// If v is less than 0, it returns 0. If v is greater than 255, it returns 255.
// Otherwise, it returns v converted to uint8.
func clamp8(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func ensureRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)
	return dst
}

func applyPaletteLuminanceFilter(img image.Image, mapper func(lum uint8) (uint8, uint8, uint8)) image.Image {
	src := ensureRGBA(img)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		srcOffset := src.PixOffset(bounds.Min.X, y)
		dstOffset := dst.PixOffset(bounds.Min.X, y)

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r := src.Pix[srcOffset]
			g := src.Pix[srcOffset+1]
			b := src.Pix[srcOffset+2]
			a := src.Pix[srcOffset+3]

			lum := luminance8(r, g, b)
			nr, ng, nb := mapper(lum)

			dst.Pix[dstOffset] = nr
			dst.Pix[dstOffset+1] = ng
			dst.Pix[dstOffset+2] = nb
			dst.Pix[dstOffset+3] = a

			srcOffset += 4
			dstOffset += 4
		}
	}

	return dst
}

func luminance8(r, g, b uint8) uint8 {
	return uint8((299*uint32(r) + 587*uint32(g) + 114*uint32(b)) / 1000)
}

func colorToNRGBA(c color.Color) color.NRGBA {
	if nrgba, ok := c.(color.NRGBA); ok {
		return nrgba
	}

	r, g, b, a := c.RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
