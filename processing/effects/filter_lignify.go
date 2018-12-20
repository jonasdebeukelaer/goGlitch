package effects

import (
	"image"
	"image/color"
	"log"
	"math"
	"sync"

	"github.com/jonasdebeukelaer/goGlitch/processing"
)

// Lignify processes the image with the effect lignify
func Lignify(p processing.Process) error {
	log.Print("Running lignify...")

	img, err := p.GetSourceImage()
	if err != nil {
		log.Println("cancelled!")
		return err
	}

	xMax := img.Bounds().Max.X
	yMax := img.Bounds().Max.Y

	upLeft := image.Point{0, 0}
	lowRight := image.Point{xMax, yMax}

	processedImg := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var wg sync.WaitGroup
	wg.Add(yMax)
	for y := 0; y < yMax; y++ {
		go func() {
			defer wg.Done()

			for x := 0; x < xMax; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				r = uint32(math.Ceil(float64(r) * 1.1))
				processedImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			}
		}()
	}

	log.Println("done")
	wg.Wait()
	return p.SetProcessedImage(processedImg)
}
