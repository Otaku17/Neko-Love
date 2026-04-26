package filters

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/chai2010/webp"
)

// Pixelate applies a pixelation effect to the given image by dividing it into blocks of fixed size
// and replacing each block with its average color. The function returns a new image with the effect applied.
//
// Parameters:
//   - img: The source image to be pixelated.
//
// Returns:
//   - image.Image: A new image with the pixelation effect applied.
func Pixelate(img image.Image) image.Image {
	return PixelateWithBlockSize(img, 6)
}

// PixelateWithBlockSize applies the pixelation effect using a configurable block size.
// A larger block size produces a stronger pixelation effect.
func PixelateWithBlockSize(img image.Image, blockSize int) image.Image {
	src := ensureRGBA(img)
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	if blockSize < 2 {
		blockSize = 2
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += blockSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += blockSize {
			var rTotal, gTotal, bTotal, aTotal int
			count := 0

			for yy := y; yy < y+blockSize && yy < bounds.Max.Y; yy++ {
				rowOffset := src.PixOffset(x, yy)
				for xx := x; xx < x+blockSize && xx < bounds.Max.X; xx++ {
					rTotal += int(src.Pix[rowOffset])
					gTotal += int(src.Pix[rowOffset+1])
					bTotal += int(src.Pix[rowOffset+2])
					aTotal += int(src.Pix[rowOffset+3])
					count++
					rowOffset += 4
				}
			}

			rAvg := uint8(rTotal / count)
			gAvg := uint8(gTotal / count)
			bAvg := uint8(bTotal / count)
			aAvg := uint8(aTotal / count)

			for yy := y; yy < y+blockSize && yy < bounds.Max.Y; yy++ {
				rowOffset := dst.PixOffset(x, yy)
				for xx := x; xx < x+blockSize && xx < bounds.Max.X; xx++ {
					dst.Pix[rowOffset] = rAvg
					dst.Pix[rowOffset+1] = gAvg
					dst.Pix[rowOffset+2] = bAvg
					dst.Pix[rowOffset+3] = aAvg
					rowOffset += 4
				}
			}
		}
	}

	return dst
}
