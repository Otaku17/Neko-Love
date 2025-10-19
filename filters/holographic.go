package filters

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// HoloFuturistic applies a futuristic hologram style to the image.
// It boosts cool tones, adds scanlines, glow-like edge lighting and a digital distortion effect.
func HoloFuturistic(img image.Image) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	bounds.Dx()
	bounds.Dy()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			// Amplify blue and cyan tones
			coolBoost := 1.3
			newR := clampF(r8*0.6, 0, 255)
			newG := clampF(g8*coolBoost, 0, 255)
			newB := clampF(b8*coolBoost*1.2, 0, 255)

			// Add scanline glow effect
			scanlineFactor := 0.95 + 0.05*math.Sin(float64(y)/2.0)
			newR *= scanlineFactor
			newG *= scanlineFactor
			newB *= scanlineFactor

			dst.Set(x, y, color.NRGBA{
				R: uint8(clampF(newR, 0, 255)),
				G: uint8(clampF(newG, 0, 255)),
				B: uint8(clampF(newB, 0, 255)),
				A: uint8(a >> 8),
			})
		}
	}

	// Optional: simulate light bleed on edges
	applyEdgeGlow(dst)

	return dst
}

func applyEdgeGlow(img *image.RGBA) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Sobel-like edge detection and glow application
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			cx := getGray(img.RGBAAt(x+1, y)) - getGray(img.RGBAAt(x-1, y))
			cy := getGray(img.RGBAAt(x, y+1)) - getGray(img.RGBAAt(x, y-1))
			edge := math.Sqrt(float64(cx*cx + cy*cy))

			if edge > 30 {
				p := img.RGBAAt(x, y)
				img.Set(x, y, color.NRGBA{
					R: clampU8(float64(p.R) + 40),
					G: clampU8(float64(p.G) + 80),
					B: clampU8(float64(p.B) + 120),
					A: p.A,
				})
			}
		}
	}
}

func getGray(c color.RGBA) int {
	return int(0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B))
}

func clampF(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampU8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
