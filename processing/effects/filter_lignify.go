package effects

import (
	"image"
	"image/color"
	"log"
	"sync"
)

// Lignify processes the image with the effect lignify
func Lignify(img image.Image) (image.Image, error) {
	log.Print("Running lignify...")

	xMax := img.Bounds().Max.X
	yMax := img.Bounds().Max.Y

	lineCount := 80
	pxSeparation := 1.0 * xMax / lineCount
	scaler := float64(1.0*pxSeparation) / (256*256 - 1)

	processedImg := image.NewRGBA(img.Bounds())

	black := color.RGBA{0, 0, 0, 255}
	processedImg, err := Fill(processedImg, black)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for y := 0; y < yMax; y++ {
		sy := y

		if y%pxSeparation == 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for x := 0; x < xMax; x++ {
					r, g, b, a := img.At(x, sy).RGBA()
					avgBrightness := (r + g + b + a) / 1.0
					vOffset := float64(avgBrightness)*scaler - float64(pxSeparation/2)

					newY, err := Wrap(sy+int(vOffset), yMax)
					if err != nil {
						log.Fatal("error wrapping y")
						newY = y
					}
					processedImg.Set(x, newY, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				}
			}()
		}
	}
	wg.Wait()

	processedImg, err = Blur(processedImg, 0.9)
	if err != nil {
		return nil, err
	}

	log.Println("done")
	return processedImg, nil
}
