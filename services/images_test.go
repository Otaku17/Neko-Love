package services

import (
	"image"
	"image/color"
	"image/gif"
	"testing"
)

func TestProcessGIFUsesFullCanvasBounds(t *testing.T) {
	t.Parallel()

	first := image.NewPaletted(image.Rect(0, 0, 2, 2), []color.Color{color.Transparent, color.White})
	second := image.NewPaletted(image.Rect(1, 1, 4, 3), []color.Color{color.Transparent, color.White})
	first.SetColorIndex(0, 0, 1)
	second.SetColorIndex(1, 1, 1)
	fullBounds := image.Rect(0, 0, 4, 3)

	src := &gif.GIF{
		Config:   image.Config{Width: 4, Height: 3},
		Image:    []*image.Paletted{first, second},
		Delay:    []int{5, 10},
		Disposal: []byte{gif.DisposalNone, gif.DisposalBackground},
	}

	result, err := ProcessGIF("negative", src, DefaultFilterOptions())
	if err != nil {
		t.Fatalf("ProcessGIF failed: %v", err)
	}

	if got := result.Image[0].Bounds(); got != fullBounds {
		t.Fatalf("expected first frame bounds %v, got %v", fullBounds, got)
	}
	if got := result.Image[1].Bounds(); got != fullBounds {
		t.Fatalf("expected second frame bounds %v, got %v", fullBounds, got)
	}
}

func TestProcessGIFCompositesPartialFramesOnFullCanvas(t *testing.T) {
	t.Parallel()

	fullBounds := image.Rect(0, 0, 3, 3)
	first := image.NewPaletted(fullBounds, []color.Color{color.Transparent, color.RGBA{255, 0, 0, 255}})
	first.SetColorIndex(0, 0, 1)

	second := image.NewPaletted(image.Rect(1, 1, 2, 2), []color.Color{color.Transparent, color.RGBA{0, 255, 0, 255}})
	second.SetColorIndex(1, 1, 1)

	src := &gif.GIF{
		Config: image.Config{Width: 3, Height: 3},
		Image:  []*image.Paletted{first, second},
		Delay:  []int{1, 1},
	}

	result, err := ProcessGIF("negative", src, DefaultFilterOptions())
	if err != nil {
		t.Fatalf("ProcessGIF failed: %v", err)
	}

	if got := result.Image[1].Bounds(); got != fullBounds {
		t.Fatalf("expected composited frame bounds %v, got %v", fullBounds, got)
	}

	r, g, b, a := result.Image[1].At(0, 0).RGBA()
	if a == 0 {
		t.Fatal("expected first frame content to remain visible on second composited frame")
	}
	if r == 0 && g == 0 && b == 0 {
		t.Fatal("expected filtered composited pixel to contain color information")
	}
}
