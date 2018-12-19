package imageprocessing

import (
	"errors"
	"image"
	_ "image/jpeg" //to make sure jpeg images can be decoded
	"image/png"
	"log"
	"os"
	"strings"
)

// ProcessImage processes an image specified by filename
func ProcessImage(filename string) error {
	fileReader, err := os.Open("uploads/" + filename)
	if err != nil {
		log.Fatalf("Error loading image for processing: %v", err)
	}
	defer fileReader.Close()

	m, _, err := image.Decode(fileReader)
	if err != nil {
		log.Fatalf("Error decoding image for processing: %v", err)
	}

	processedImg, err := lignify(m)
	if err != nil {
		return err
	}

	err = saveImage(processedImg, filename)
	if err != nil {
		return errors.New("Error saving image: " + err.Error())
	}

	log.Println("completing")
	return nil
}

func saveImage(img image.Image, filename string) error {
	log.Printf("Saving processed image %s", filename)
	f, err := os.Create("processed_images/" + strings.Split(filename, ".")[0] + "_processed.png")
	if err != nil {
		return err
	}
	png.Encode(f, img)

	return nil
}
