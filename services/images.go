package services

import (
	"errors"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"neko-love/filters"

	"github.com/chai2010/webp"
	"github.com/gofiber/fiber/v2"
)

// ProcessGIF applies a specified filter to each frame of a given GIF image.
// It processes each frame by converting it to RGBA, applying the filter, and then
// converting it back to a paletted image while preserving transparency. The function
// also maintains the original GIF's loop count, frame delays, and disposal methods.
//
// Parameters:
//   - filterName: The name of the filter to apply to each frame.
//   - g: Pointer to the gif.GIF object to be processed.
//
// Returns:
//   - A pointer to a new gif.GIF object with the filter applied to each frame.
//   - An error if the input GIF has no frames or if processing fails.
func ProcessGIF(filterName string, g *gif.GIF) (*gif.GIF, error) {
	if len(g.Image) == 0 {
		return nil, errors.New("GIF has no frames")
	}

	result := &gif.GIF{
		LoopCount: g.LoopCount,
		Image:     make([]*image.Paletted, 0, len(g.Image)),
		Delay:     make([]int, 0, len(g.Delay)),
		Disposal:  make([]byte, 0, len(g.Disposal)),
	}

	for i, frame := range g.Image {
		bounds := g.Image[0].Bounds()
		rgba := image.NewRGBA(bounds)

		draw.Draw(rgba, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

		filtered := ApplyFilter(filterName, rgba)

		palettedFrame := rgbaToPalettedWithTransparency(filtered)

		palettedFrame.Rect = bounds

		result.Image = append(result.Image, palettedFrame)

		if i < len(g.Delay) {
			result.Delay = append(result.Delay, g.Delay[i])
		} else {
			result.Delay = append(result.Delay, 0)
		}
		if i < len(g.Disposal) {
			result.Disposal = append(result.Disposal, g.Disposal[i])
		} else {
			result.Disposal = append(result.Disposal, gif.DisposalNone)
		}
	}

	return result, nil
}

// ApplyFilter applies the specified filter to the provided image and returns the resulting image.
// Supported filters include: "blurple", "fuchsia", "glitch", "neon", "deepfry", "posterize", "pixelate", "vaporwave", "anime_outline".... .
// If an unknown filter is provided, the original image is returned unmodified.
//
// Parameters:
//   - filter: the name of the filter to apply.
//   - img: the image.Image to which the filter will be applied.
//
// Returns:
//   - image.Image: the filtered image.
func ApplyFilter(filter string, img image.Image) image.Image {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	switch filter {
	case "blurple":
		return filters.Blurple(rgba)
	case "fuchsia":
		return filters.Fuchsia(rgba)
	case "glitch":
		return filters.Glitch(rgba)
	case "poppink":
		return filters.PopPink(rgba)
	case "deepfry":
		return filters.Deepfry(rgba)
	case "posterize":
		return filters.Posterize(rgba)
	case "pixelate":
		return filters.Pixelate(rgba)
	case "vaporwave":
		return filters.Vaporwave(rgba)
	case "anime_outline":
		return filters.AnimeOutline(rgba)
	case "crimson":
		return filters.Crimson(rgba)
	case "amber":
		return filters.Amber(rgba)
	case "mint":
		return filters.Mint(rgba)
	case "aqua":
		return filters.Aqua(rgba)
	case "sunset":
		return filters.Sunset(rgba)
	case "bubblegum":
		return filters.Bubblegum(rgba)
	case "negative":
		return filters.Negative(rgba)
	case "greyscale":
		return filters.Greyscale(rgba)
	case "holographic":
		return filters.HoloFuturistic(rgba)
	default:
		return rgba
	}
}

// rgbaToPalettedWithTransparency converts an RGBA image to a paletted image using the Plan9 palette,
// ensuring that fully or partially transparent pixels are mapped to the first palette entry (index 0),
// which is set to fully transparent. The function applies Floyd-Steinberg dithering for color quantization.
// It returns the resulting *image.Paletted.
func rgbaToPalettedWithTransparency(img image.Image) *image.Paletted {
	bounds := img.Bounds()

	p := make(color.Palette, len(palette.Plan9))
	copy(p, palette.Plan9)

	p[0] = color.RGBA{0, 0, 0, 0}

	palettedImg := image.NewPaletted(bounds, p)

	draw.FloydSteinberg.Draw(palettedImg, bounds, img, image.Point{})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if uint8(a>>8) < 255 {
				palettedImg.SetColorIndex(x, y, 0)
			}
		}
	}

	return palettedImg
}

// EncodeAndSetContentType encodes the provided image.Image into the specified format
// ("jpeg", "png", or "webp") and writes it to the Fiber context response body,
// setting the appropriate Content-Type header. If the format is unrecognized,
// it defaults to encoding as PNG. Returns an error if encoding fails.
func EncodeAndSetContentType(c *fiber.Ctx, img image.Image, formatStr string) error {
	switch formatStr {
	case "jpeg":
		c.Set("Content-Type", "image/jpeg")
		return jpeg.Encode(c.Context().Response.BodyWriter(), img, &jpeg.Options{Quality: 90})
	case "png":
		c.Set("Content-Type", "image/png")
		return png.Encode(c.Context().Response.BodyWriter(), img)
	case "webp":
		c.Set("Content-Type", "image/webp")
		return webp.Encode(c.Context().Response.BodyWriter(), img, &webp.Options{Lossless: true})
	default:
		// Fallback: PNG
		c.Set("Content-Type", "image/png")
		return png.Encode(c.Context().Response.BodyWriter(), img)
	}
}
