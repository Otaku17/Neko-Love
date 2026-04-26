package filters

import (
	"image"
	"image/color"
	"testing"
)

func TestNegativePreservesAlphaAndInvertsRGB(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.SetRGBA(0, 0, color.RGBA{R: 10, G: 20, B: 30, A: 40})

	got := colorToNRGBA(Negative(src).At(0, 0))
	want := color.NRGBA{R: 245, G: 235, B: 225, A: 40}

	if got != want {
		t.Fatalf("unexpected pixel: got %#v want %#v", got, want)
	}
}

func TestGreyscaleProducesEqualChannels(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.SetRGBA(0, 0, color.RGBA{R: 100, G: 150, B: 200, A: 255})

	got := colorToNRGBA(Greyscale(src).At(0, 0))
	if got.R != got.G || got.G != got.B {
		t.Fatalf("expected greyscale channels to match, got %#v", got)
	}
	if got.A != 255 {
		t.Fatalf("expected alpha to be preserved, got %#v", got)
	}
}

func TestPosterizeQuantizesChannels(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.SetRGBA(0, 0, color.RGBA{R: 130, G: 65, B: 255, A: 200})

	got := colorToNRGBA(Posterize(src).At(0, 0))
	want := color.NRGBA{R: 128, G: 64, B: 192, A: 200}

	if got != want {
		t.Fatalf("unexpected posterized pixel: got %#v want %#v", got, want)
	}
}

func TestAmberDoesNotCollapseBrightImageToDark(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.SetRGBA(0, 0, color.RGBA{R: 255, G: 220, B: 180, A: 255})

	got := colorToNRGBA(Amber(src).At(0, 0))
	if got.R == 35 && got.G == 39 && got.B == 42 {
		t.Fatalf("amber filter unexpectedly collapsed bright pixel to dark tone: %#v", got)
	}
}
