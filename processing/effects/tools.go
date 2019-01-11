package effects

import (
	"errors"
	"image"
	"image/color"
	"sync"

	"github.com/disintegration/gift"
)

var (
	// ErrMaxCannotBeZero returned when max is zero
	ErrMaxCannotBeZero = errors.New("Max cannot be zero")
	// ErrImageSizeIsZero returned when image size is 0
	ErrImageSizeIsZero = errors.New("Image cannot have size 0")
)

// Wrap wraps around the edger of an image and back around
// the other side
func Wrap(a, aMax int) (int, error) {
	if aMax < 1 {
		return 0, ErrMaxCannotBeZero
	}
	for a < 0 || a > aMax {
		if a < 0 {
			a = aMax + a
		} else if a > aMax {
			a = a - aMax
		}
	}
	return a, nil
}

// Fill fills an input image with the specificed colour
func Fill(img *image.RGBA, colour color.RGBA) (*image.RGBA, error) {
	if img.Bounds() == image.Rect(0, 0, 0, 0) {
		return nil, ErrImageSizeIsZero
	}

	xMax := img.Bounds().Max.X
	yMax := img.Bounds().Max.Y

	filledImg := image.NewRGBA(img.Bounds())

	var wg sync.WaitGroup
	wg.Add(yMax)
	for y := 0; y < yMax; y++ {
		sy := y
		go func() {
			defer wg.Done()

			for x := 0; x < xMax; x++ {
				filledImg.Set(x, sy, colour)
			}
		}()
	}
	wg.Wait()
	return filledImg, nil
}

// Blur uses a Gaussian blur on the input image, given a sigma.
// sigma defines how much to blur the image
func Blur(src image.Image, sigma float32) (*image.RGBA, error) {
	if src.Bounds() == image.Rect(0, 0, 0, 0) {
		return nil, ErrImageSizeIsZero
	}

	g := gift.New(
		gift.GaussianBlur(sigma),
	)

	dst := image.NewRGBA(g.Bounds(src.Bounds()))

	g.Draw(dst, src)
	return dst, nil
}
