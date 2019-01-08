package processing

import (
	"testing"
)

func Test_EffectLignify(t *testing.T) {
	effect := "lignify"
	p, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image processor: %v", err)
	}

	err = p.SetEffect(EffectLignify)
	if err != nil {
		t.Fatalf("Could not set the effect '%s': %v", effect, err)
	}

	err = p.ProcessImage()
	if err != nil {
		t.Fatalf("Could not process image using effect '%s': %v", effect, err)
	}
}

func Test_EffectRandomMuddle(t *testing.T) {
	effect := "randomMuddle"
	p, err := New(testImageFilename)
	if err != nil {
		t.Fatalf("Error create new image processor: %v", err)
	}

	err = p.SetEffect(EffectRandomMuddle)
	if err != nil {
		t.Fatalf("Could not set the effect '%s': %v", effect, err)
	}

	err = p.ProcessImage()
	if err != nil {
		t.Fatalf("Could not process image using effect '%s': %v", effect, err)
	}
}
