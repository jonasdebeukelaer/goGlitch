package processing

import (
	"testing"
)

func Test_scaleImage(t *testing.T) {
	img := getImage(testImageFilename)

	scaledImg, err := scaleImage(img, 400)
	if err != nil {
		t.Fatalf("Error rescaling image: %v", err)
	}

	newXMax := scaledImg.Bounds().Max.X
	newYMax := scaledImg.Bounds().Max.Y

	if newXMax > 400 || newYMax > 400 {
		t.Fatal("Image not scaled according to max dimention")
	}

	if newXMax != 400 && newYMax != 400 {
		t.Fatal("One of the dimentions not rescaled to 400")
	}
}

func Test_scaleImageToZeroShouldFail(t *testing.T) {
	img := getImage(testImageFilename)

	_, err := scaleImage(img, 0)
	if err == nil {
		t.Fatal("Able to try rescale to 0")
	}
}
