package processor

import (
	"errors"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
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

// New creates an instance of an image manupulation process
func New(filename string) (processing.Process, error) {
	_, realFilename := filepath.Split(filename)
	filenameParts := strings.Split(realFilename, ".")
	imgName := strings.Join(filenameParts[:len(filenameParts)-1], ".")

	p := &processor{
		sourceImgFilename: filename,
		sourceImgName:     imgName,
	}

	err := p.setSourceImage(filename)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *processor) setSourceImage(filename string) error {
	fileReader, err := os.Open(filename)
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
	// e, ok := processing.EffectLignify
	// if !ok {
	// 	return errors.New("effect '" + effect + "' not valid")
	// }
	p.effect = effect
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
	targetDirectory := "storage/processed_images/"
	targetFilepath := targetDirectory + p.processedImgFilename

	err := createDirectoryIfNotExists(targetDirectory)
	if err != nil {
		return errors.New("Error creating target directory '" + targetDirectory + "'")
	}

	f, err := os.Create(targetFilepath)
	if err != nil {
		return err
	}
	return png.Encode(f, p.processedImg)
}

func createDirectoryIfNotExists(targetDirectory string) error {
	_, err := os.Stat(targetDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(targetDirectory, 0755)
		}
		return err
	}
	return nil
}
