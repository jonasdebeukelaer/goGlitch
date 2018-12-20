package processing

import (
	"image"
	_ "image/jpeg" //to make sure jpeg images can be decoded
)

// Process represents the steps required to apply an effect or
// filter on a source image to create a new processed image
type Process interface {
	GetSourceImage() (image.Image, error)

	SetEffect(effect string) error

	ProcessImage() error

	SetProcessedImage(image.Image) error
	GetProcessedImage() (image.Image, error)
	GetProcessedImageFilename() (string, error)
}