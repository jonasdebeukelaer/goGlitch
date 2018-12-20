package processor

import (
	"errors"
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/jonasdebeukelaer/goGlitch/processing"
	"github.com/jonasdebeukelaer/goGlitch/processing/effects"
)

type processor struct {
	sourceImgFilename string
	sourceImgName     string
	sourceImgFiletype string
	sourceImg         image.Image

	effect             string
	processingComplete bool

	processedImgFilename string
	processedImg         image.Image
}

// New creates an instance of an image manupulation process Process
func New(filename string) (processing.Process, error) {
	filenameParts := strings.Split(filename, ".")
	imgName := strings.Join(filenameParts[:len(filenameParts)-1], ".")

	p := &processor{
		sourceImgFilename: filename,
		sourceImgName:     imgName,
	}

	p.setSourceImage(filename)

	return p, nil
}

func (p *processor) setSourceImage(filename string) error {
	fileReader, err := os.Open("uploads/" + filename)
	if err != nil {
		return errors.New("Error loading image for processing: " + err.Error())
	}
	defer fileReader.Close()

	img, filetype, err := image.Decode(fileReader)
	if err != nil {
		return errors.New("Error decoding image for processing: " + err.Error())
	}

	p.sourceImg = img
	p.sourceImgFiletype = filetype
	return nil
}

func (p *processor) SetEffect(effect string) error {
	e, ok := processing.Options[effect]
	if !ok {
		return errors.New("effect '" + effect + "' not valid")
	}
	p.effect = e
	return nil
}

func (p *processor) ProcessImage() error {
	if p.processingComplete {
		return errors.New("Processing already complete")
	}
	if p.sourceImg == nil {
		return errors.New("Source image not set")
	}

	err := effects.Lignify(p) // TODO change this once more effects are added!
	if err != nil {
		return err
	}
	p.processingComplete = true
	return nil
}

func (p processor) GetSourceImage() (image.Image, error) {
	if p.sourceImg == nil {
		return nil, errors.New("Source image not yet set")
	}
	return p.sourceImg, nil
}

func (p *processor) SetProcessedImage(img image.Image) error {
	p.processedImg = img
	return p.saveImage()
}

func (p processor) GetProcessedImage() (image.Image, error) {
	if !p.processingComplete {
		return nil, errors.New("Processing has not yet occured")
	}
	return p.processedImg, nil
}

func (p processor) GetProcessedImageFilename() (string, error) {
	if !p.processingComplete {
		return "", errors.New("Processing has not yet occured")
	}
	return p.processedImgFilename, nil
}

func (p *processor) saveImage() error {
	log.Printf("Saving processed image %s", p.sourceImgFilename)
	p.processedImgFilename = p.sourceImgName + "_processed.png"

	f, err := os.Create("processed_images/" + p.processedImgFilename)
	if err != nil {
		return err
	}
	return png.Encode(f, p.processedImg)
}