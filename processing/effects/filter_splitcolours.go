package effects

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strconv"
	"sync"
)

// SplitColours processes the image with the effect split colours
func SplitColours(img image.Image, params map[string]string) (image.Image, error) {
	log.Print("Running split colours...")

	offsetPercent := 5.0
	if params != nil {
		if op, ok := params["offsetPercent"]; ok {
			offperc, err := strconv.ParseFloat(op, 64)
			if err != nil {
				return nil, fmt.Errorf("could not extract param: %v", err)
			}
			offsetPercent = offperc
			log.Printf("...with offsetPercent %v", offsetPercent)
		}
	}

	xMax := img.Bounds().Max.X
	yMax := img.Bounds().Max.Y

	pxOffset := offsetPercent / 100.0 * float64(xMax)

	processedImg := image.NewRGBA(img.Bounds())

	black := color.RGBA{0, 0, 0, 255}
	processedImg, err := Fill(processedImg, black)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for y := 0; y < yMax; y++ {
		sy := y
		wg.Add(1)
		go func() {
			defer wg.Done()

			for x := 0; x < xMax; x++ {
				rx := Wrap(x-int(pxOffset), xMax)
				rsy := Wrap(sy-int(pxOffset), yMax)

				bx := Wrap(x+int(pxOffset), xMax)
				bsy := Wrap(sy+int(pxOffset), yMax)

				r, _, _, _ := img.At(rx, rsy).RGBA()
				_, g, _, a := img.At(x, sy).RGBA()
				_, _, b, _ := img.At(bx, bsy).RGBA()

				processedImg.Set(x, sy, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			}
		}()
	}
	wg.Wait()

	return processedImg, nil
}
