package processing

import (
	"errors"
	"image"
	_ "image/jpeg" // required for decoding
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Processor struct {
	sourceImgFilename string
	sourceImgName     string
	sourceImgFiletype string
	sourceImg         image.Image

	processingComplete   bool
	processedImgFilename string
	processedImg         image.Image
}

type EffectLayer struct {
	Key    string
	Params map[string]string
}

// New creates an instance of an image manupulation process
func New(filename string) (*Processor, error) {
	_, realFilename := filepath.Split(filename)
	filenameParts := strings.Split(realFilename, ".")
	imgName := strings.Join(filenameParts[:len(filenameParts)-1], ".")

	p := &Processor{
		sourceImgFilename: filename,
		sourceImgName:     imgName,
	}

	err := p.setSourceImage(filename)
	if err != nil {
		return &Processor{}, err
	}

	return p, nil
}

func (p *Processor) setSourceImage(filename string) error {
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

func (p *Processor) ProcessImage(layers []*EffectLayer) error {
	if p.processingComplete {
		return errors.New("Processing already complete")
	}
	if p.sourceImg == nil {
		return errors.New("Source image not set")
	}

	scaledImage, err := scaleImage(p.sourceImg, 1000)
	if err != nil {
		return err
	}

	filenameSuffix := ""
	workingImg := scaledImage
	for _, layer := range layers {
		filenameSuffix += "_" + layer.Key
		workingImg, err = EffectMap[layer.Key](workingImg, layer.Params)
		if err != nil {
			return err
		}
	}
	p.SetProcessedImage(workingImg, filenameSuffix)

	p.processingComplete = true
	return nil
}

func (p Processor) GetSourceImage() (image.Image, error) {
	if p.sourceImg == nil {
		return nil, errors.New("Source image not yet set")
	}
	return p.sourceImg, nil
}

func (p *Processor) SetProcessedImage(img image.Image, filenameSuffix string) error {
	p.processedImg = img
	return p.saveImage(filenameSuffix)
}

func (p Processor) GetProcessedImage() (image.Image, error) {
	if !p.processingComplete {
		return nil, errors.New("Processing has not yet occured")
	}
	return p.processedImg, nil
}

func (p Processor) GetProcessedImageFilename() (string, error) {
	if !p.processingComplete {
		return "", errors.New("Processing has not yet occured")
	}
	return p.processedImgFilename, nil
}

func (p *Processor) saveImage(filenameSuffix string) error {
	p.processedImgFilename = p.sourceImgName + "_" + filenameSuffix + ".png"
	log.Printf("Saving processed image %s", p.processedImgFilename)

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
