package processing

import (
	"errors"
	"image"

	"github.com/disintegration/gift"
)

func scaleImage(img image.Image, maxDim int) (image.Image, error) {
	if maxDim < 1 {
		return nil, errors.New("Cannot rescale to less than 1")
	}

	xMax, yMax := float32(img.Bounds().Max.X), float32(img.Bounds().Max.Y)

	var newX, newY int
	if xMax > yMax {
		newX = int(maxDim)
		newY = int(float32(maxDim) * (yMax / xMax))
	} else {
		newX = int(float32(maxDim) * (xMax / yMax))
		newY = int(maxDim)
	}

	g := gift.New(
		gift.Resize(newX, newY, gift.LanczosResampling),
	)

	scaledImg := image.NewRGBA(image.Rect(0, 0, newX, newY))
	g.Draw(scaledImg, img)

	return scaledImg, nil
}
