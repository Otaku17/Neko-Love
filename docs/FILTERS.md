# Filter Development Guide

This project keeps filters simple on purpose. A good filter should be easy to read, cheap to run, and predictable on PNG, JPEG, WEBP, and GIF frame processing.

## Recommended structure

1. Add one file per filter in `filters/`.
2. Export a single function like `func MyFilter(img image.Image) image.Image`.
3. Return a new image instead of mutating shared input state.
4. Register the filter in `services.ApplyFilter` inside [services/images.go](/c:/Users/steve/Desktop/plugin/vue-theme-plugin/Neko-Love/services/images.go:71).

## Good practices

- Prefer a single full-image pass when possible.
- Work with RGBA-compatible output so alpha stays predictable.
- Preserve transparency unless the effect explicitly replaces it.
- Avoid random effects unless the style really needs them.
- Keep allocations low: do not create temporary images inside the inner pixel loop.
- If the filter is expensive, keep math simple and clamp values carefully.
- Make the result deterministic for most filters so tests and previews stay stable.

## Suggested template

```go
package filters

import (
    "image"
    "image/color"
)

func MyFilter(img image.Image) image.Image {
    bounds := img.Bounds()
    dst := image.NewRGBA(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()

            dst.Set(x, y, color.NRGBA{
                R: uint8(r >> 8),
                G: uint8(g >> 8),
                B: uint8(b >> 8),
                A: uint8(a >> 8),
            })
        }
    }
V
    return dst
}
```

## Testing checklist

- Verify the filter works on a tiny PNG fixture.
- Verify transparency is still correct if the source contains alpha.
- If the filter is used for GIFs, make sure it behaves frame-by-frame without assuming all frames have the same bounds.
- If the filter uses randomness, consider injecting or controlling the random source for tests.

## When not to do it this way

- If several filters share the same heavy logic, move that logic into a helper.
- If a filter needs parameters, consider a separate endpoint shape later instead of hardcoding one-off query parsing into every filter.
