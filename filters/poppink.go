package filters

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"

	_ "github.com/chai2010/webp"
)

// edgeDetect calculates the average color difference between the pixel at (x, y)
// and its immediate neighbors (up, down, left, right) in the given image.
// It returns a float64 representing the average edge strength at that pixel.
// If the pixel has no valid neighbors (e.g., at the image edge), it returns 0.
func edgeDetect(img image.Image, x, y int) float64 {
	bounds := img.Bounds()
	c := img.At(x, y)
	r0, g0, b0, _ := c.RGBA()
	r0f := float64(r0 >> 8)
	g0f := float64(g0 >> 8)
	b0f := float64(b0 >> 8)

	var diffSum float64
	var count int

	deltas := []image.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for _, d := range deltas {
		nx, ny := x+d.X, y+d.Y
		if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
			nc := img.At(nx, ny)
			r1, g1, b1, _ := nc.RGBA()
			r1f := float64(r1 >> 8)
			g1f := float64(g1 >> 8)
			b1f := float64(b1 >> 8)

			diff := math.Abs(r0f-r1f) + math.Abs(g0f-g1f) + math.Abs(b0f-b1f)
			diffSum += diff / 3.0
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return diffSum / float64(count)
}

// PopPink applies a vibrant "pop pink" filter effect to the given image.
// The effect consists of two main components:
//  1. A neon blue color boost applied to the image's base colors, increasing
//     the intensity of blue and red channels for a vivid look.
//  2. A neon red halo effect around detected edges, giving the image a glowing,
//     stylized outline. The halo is created by edge detection, thresholding,
//     and dilation to expand the glow around edges.
//
// The resulting image combines the enhanced base colors with the glowing edge
// halo, producing a visually striking, pop-art inspired effect.
//
// Parameters:
//
//	img image.Image - The source image to apply the filter to.
//
// Returns:
//
//	image.Image - The filtered image with the pop pink effect applied.
func PopPink(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	neonBlue := color.NRGBA{R: 80, G: 180, B: 255, A: 0}
	neonRed := color.NRGBA{R: 255, G: 40, B: 60, A: 0}

	halo := image.NewRGBA(bounds)

	edgeThreshold := 20.0
	outerAlpha := uint8(60)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			edgeVal := edgeDetect(img, x, y)

			if edgeVal > edgeThreshold {
				halo.Set(x, y, color.NRGBA{neonRed.R, neonRed.G, neonRed.B, outerAlpha})
			} else {
				halo.Set(x, y, color.NRGBA{0, 0, 0, 0})
			}
		}
	}

	dilateOffsets := []image.Point{
		{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {-1, -1}, {1, -1}, {-1, 1},
	}
	haloDilated := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			maxA := uint8(0)
			var r, g, b uint8
			for _, off := range dilateOffsets {
				nx, ny := x+off.X, y+off.Y
				if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
					cr, cg, cb, ca := halo.At(nx, ny).RGBA()
					a8 := uint8(ca >> 8)
					if a8 > maxA {
						maxA = a8
						r = uint8(cr >> 8)
						g = uint8(cg >> 8)
						b = uint8(cb >> 8)
					}
				}
			}
			haloDilated.Set(x, y, color.NRGBA{r, g, b, maxA})
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			origR, origG, origB, a := img.At(x, y).RGBA()
			a8 := uint8(a >> 8)

			r := clamp8(int(float64(origR>>8)*1.5) + int(neonBlue.R/2))
			g := clamp8(int(float64(origG>>8)*1.2) + int(neonBlue.G/2))
			b := clamp8(int(float64(origB>>8)*1.8) + int(neonBlue.B/2))

			hr, hg, hb, ha := haloDilated.At(x, y).RGBA()
			haF := float64(ha>>8) / 255.0

			finalR := uint8(float64(r)*(1-haF) + float64(hr>>8)*haF)
			finalG := uint8(float64(g)*(1-haF) + float64(hg>>8)*haF)
			finalB := uint8(float64(b)*(1-haF) + float64(hb>>8)*haF)

			dst.Set(x, y, color.NRGBA{finalR, finalG, finalB, a8})
		}
	}

	return dst
}
