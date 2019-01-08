package effects

import (
	"image"
	"log"
	"math"

	"github.com/disintegration/gift"
)

// Lignify processes the image with the effect lignify
func RandomMuddle(img image.Image) (image.Image, error) {
	log.Print("Running random muddle...")

	g := gift.New(
		gift.ColorFunc(randomColourFunc),
	)

	processedImg := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(processedImg, img)
	return processedImg, nil
}

func randomColourFunc(r0, g0, b0, a0 float32) (r, g, b, a float32) {
	return r0 * 2.0, float32(math.Pow(float64(g0), 1.1)), g0 - b0, a0
}
