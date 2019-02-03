package effects

import (
	"fmt"
	"image"
	"log"
	"sync"
)

// Threedee splits the image into two layers and an effect is applied to the top layer only
// TODO:
// 1. allow any effect to be applied to one layer
// 2. make threshold adjustable
// 3. make sensible threshold calculated automatically
// 4. find more advanced techniques for selecting layer e.g. some kind of gaussian blurring or edge detection
func Threedee(img image.Image, param map[string]string) (image.Image, error) {
	log.Print("Running threedee...")

	xMax := img.Bounds().Max.X
	yMax := img.Bounds().Max.Y

	threshold := map[string]uint32{
		"r":   220,
		"g":   220,
		"b":   220,
		"all": 80000,
	}
	//dLen := 10

	layerSplit := splitByTheshold(img, xMax, yMax, threshold)
	layerSplit = tidyThreshold(img, layerSplit)

	processedImg, err := Lignify(img, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("error processing top layer: %v", err)
	}

	processedImg = combineLayers(img, processedImg, xMax, yMax, layerSplit)

	return processedImg, nil
}

func combineLayers(img image.Image, processedImg image.Image, xMax, yMax int, layerSplit [][]uint8) image.Image {
	combinedImg := image.NewRGBA(img.Bounds())

	var wg sync.WaitGroup
	for y := 0; y < yMax; y++ {
		sy := y
		wg.Add(1)
		go func() {
			defer wg.Done()

			for x := 0; x < xMax; x++ {

				if layerSplit[sy][x] == 1 {
					combinedImg.Set(x, sy, processedImg.At(x, sy))
				} else {
					combinedImg.Set(x, sy, img.At(x, sy))
				}
			}
		}()
	}
	wg.Wait()

	return combinedImg
}

func splitByTheshold(img image.Image, xMax, yMax int, threshold map[string]uint32) [][]uint8 {
	layerSplit := make([][]uint8, yMax)
	var wg sync.WaitGroup
	for y := 0; y < yMax; y++ {
		sy := y
		wg.Add(1)
		go func() {
			defer wg.Done()

			layerSplitRow := make([]uint8, xMax)
			for x := 0; x < xMax; x++ {

				r, g, b, _ := img.At(x, sy).RGBA()

				if r+g+b > threshold["all"] {
					layerSplitRow[x] = 1
				} else {
					layerSplitRow[x] = 0
				}
			}
			layerSplit[sy] = layerSplitRow
		}()
	}
	wg.Wait()

	return layerSplit
}

func tidyThreshold(img image.Image, layerSplit [][]uint8) [][]uint8 {
	sLayerSplit := layerSplit

	var wg sync.WaitGroup
	for y := 1; y < len(layerSplit)-1; y++ {
		sy := y
		wg.Add(1)
		go func() {
			defer wg.Done()

			for x := 1; x < len(layerSplit[0])-1; x++ {

				kernel := make([][]uint8, 3)
				kernel[0] = sLayerSplit[sy-1][x-1 : x+2]
				kernel[1] = sLayerSplit[sy][x-1 : x+2]
				kernel[2] = sLayerSplit[sy+1][x-1 : x+2]

				pixel := kernel[1][1]
				n := kernel[0][1] + kernel[1][0] + kernel[2][1] + kernel[1][2]

				if pixel == 0 && n >= 3 {
					layerSplit[sy][x] = 1
				} else if pixel == 1 && n <= 1 {
					layerSplit[sy][x] = 0
				}

			}
		}()
	}
	wg.Wait()

	return layerSplit
}
