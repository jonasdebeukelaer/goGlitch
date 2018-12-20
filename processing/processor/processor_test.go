package processor

import (
	"image"
	"os"
	"testing"

	"github.com/jonasdebeukelaer/goGlitch/processing"
)

func init() {
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}
}

var testImageFilename = "storage/test_uploads/test_upload.jpg"

func getImage(filename string) image.Image {
	fr, _ := os.Open(testImageFilename)
	img, _, _ := image.Decode(fr)
	return img
}
func Test_initImageProcess(t *testing.T) {
	_, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image process: %v", err)
	}
}

func Test_ErrorOnNewIfImageNotExist(t *testing.T) {
	testNoneExistantImageFilename := "storage/doesnt/exist.here"
	_, err := New(testNoneExistantImageFilename)
	if err == nil {
		t.Fatal("No error returned trying to load non-existant image")
	}
}

func Test_GetSourceImg(t *testing.T) {
	p, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image process: %v", err)
	}

	sourceImg, err := p.GetSourceImage()
	if err != nil {
		t.Fatalf("Error retrieving source image: %v", err)
	}

	img := getImage(testImageFilename)
	if sourceImg.Bounds() != img.Bounds() {
		t.Fatalf("Expected image bounds not equal: %+v, %+v", sourceImg.Bounds(), img.Bounds())
	}
}

func Test_SetEffect(t *testing.T) {
	p, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image process: %v", err)
	}

	err = p.SetEffect(processing.EffectLignify)
	if err != nil {
		t.Fatalf("Could not set the effect: %v", err)
	}
}

func Test_ProcessAndRetrieveImage(t *testing.T) {
	p, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image process: %v", err)
	}

	err = p.SetEffect(processing.EffectLignify)
	if err != nil {
		t.Fatalf("Could not set the effect: %v", err)
	}

	err = p.ProcessImage()
	if err != nil {
		t.Fatalf("Could not process image: %v", err)
	}

	_, err = p.GetProcessedImage()
	if err != nil {
		t.Fatalf("Could not get processed image: %v", err)
	}

	_, err = p.GetProcessedImageFilename()
	if err != nil {
		t.Fatalf("Could not get processed image filename: %v", err)
	}
}

func Test_CantProcessimageTwice(t *testing.T) {
	p, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image process: %v", err)
	}

	err = p.SetEffect(processing.EffectLignify)
	if err != nil {
		t.Fatalf("Could not set the effect: %v", err)
	}

	err = p.ProcessImage()
	if err != nil {
		t.Fatalf("Could not process image: %v", err)
	}

	err = p.ProcessImage()
	if err == nil {
		t.Fatalf("Able to run processing twice")
	}

}

func Test_CantRetrieveProcessedImageBeforeProcessing(t *testing.T) {
	p, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image process: %v", err)
	}

	_, err = p.GetProcessedImage()
	if err == nil {
		t.Fatal("Able to retrieve processed image before it was processed")
	}

	_, err = p.GetProcessedImageFilename()
	if err == nil {
		t.Fatal("Able to retrieve processed image filename before it was processed")
	}
}
