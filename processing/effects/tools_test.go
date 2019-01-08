package effects

import (
	"fmt"
	"image"
	"image/color"
	"testing"
)

func TestWrapping(t *testing.T) {
	testCases := []struct {
		ty          int
		yMax        int
		expectedOut int
		expectedErr error
	}{
		{0, 0, 0, ErrMaxCannotBeZero},
		{10, 10, 10, nil},
		{1, 0, 0, ErrMaxCannotBeZero},
		{95, 20, 15, nil},
		{-20, 100, 80, nil},
		{20, 19, 1, nil},
	}

	for _, tc := range testCases {
		tc := tc
		msg := fmt.Sprintf("y=%v yMax=%v expect %v and expect err %v", tc.ty, tc.yMax, tc.expectedOut, tc.expectedErr)
		t.Run(msg, func(tt *testing.T) {
			tt.Parallel()

			actualOut, ActualErr := Wrap(tc.ty, tc.yMax)

			if ActualErr != tc.expectedErr {
				tt.Errorf("wanted terror '%v', got error '%v'", tc.expectedErr, ActualErr)
			}

			if actualOut != tc.expectedOut {
				tt.Errorf("got %v, wanted %v", actualOut, tc.expectedOut)
			}

		})
	}
}

func Test_FillBlack(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	black := color.RGBA{0, 0, 0, 255}
	filledImg, err := Fill(img, black)
	if err != nil {
		t.Error(err)
	}

	actualColour := filledImg.At(5, 5)
	if actualColour != black {
		t.Errorf("expected colour '%v', got '%v'", black, actualColour)
	}
}

func Test_EmptyImagesReturnError(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	_, err := Fill(img, color.RGBA{1, 2, 3, 255})
	if err == nil {
		t.Errorf("No error on empty image")
	}

	_, err = Blur(img, 1.0)
	if err == nil {
		t.Errorf("No error on empty image")
	}

}

func Test_Blur(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	blurredImg, err := Blur(img, 1.0)
	if err != nil {
		t.Errorf("Error on blur: %v", err)
	}

	if blurredImg.Bounds() != image.Rect(0, 0, 10, 10) {
		t.Errorf("blurredImg not the same size as input image")
	}
}
