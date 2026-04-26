package filters

import (
	"image"
	"image/color"
	"testing"
)

func BenchmarkFilters(b *testing.B) {
	img := benchmarkImage(512, 512)

	benchmarks := []struct {
		name string
		fn   func(image.Image) image.Image
	}{
		{name: "Negative", fn: Negative},
		{name: "Greyscale", fn: Greyscale},
		{name: "Posterize", fn: Posterize},
		{name: "Blurple", fn: Blurple},
		{name: "Amber", fn: Amber},
		{name: "Aqua", fn: Aqua},
		{name: "Fuchsia", fn: Fuchsia},
		{name: "Deepfry", fn: Deepfry},
		{name: "Pixelate", fn: Pixelate},
		{name: "Glitch", fn: Glitch},
	}

	for _, bm := range benchmarks {
		bm := bm
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bm.fn(img)
			}
		})
	}
}

func benchmarkImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8((x * 13) % 256),
				G: uint8((y * 7) % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return img
}
