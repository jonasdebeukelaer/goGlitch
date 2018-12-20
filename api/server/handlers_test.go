package server

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"
)

func init() {
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}
}

func TestDefaultHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(defaultHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestCheckFileType(t *testing.T) {
	testImageType := make([]string, 1)
	testImageType[0] = "image/png"

	testMimeHeader := textproto.MIMEHeader{
		"Content-Type": testImageType,
	}

	testHandle := multipart.FileHeader{
		Filename: "test",
		Header:   testMimeHeader,
	}

	err := checkFileType(&testHandle)
	if err != nil {
		t.Errorf("Did not access filetype %s: %v", testImageType[0], err)
	}

}
