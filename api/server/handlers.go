package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/jonasdebeukelaer/goGlitch/processing"
)

type mainPageVariables struct {
}

type workPageVariables struct {
	Filename string
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/main.html")
}

func imageUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	image, handle, err := r.FormFile("image")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer image.Close()
	fmt.Print("uploading image ", handle.Filename, "...")

	err = checkFileType(handle)
	if err != nil {
		fmt.Println("error!\nImage could not be uploaded: %", err)
		return
	}

	saveImage(w, image, handle)
	fmt.Println("success!")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"filename\": \"" + handle.Filename + "\"}"))
}

func checkFileType(handle *multipart.FileHeader) error {
	mimeType := handle.Header.Get("Content-Type")

	switch mimeType {
	case "image/jpeg", "image/png":
		return nil
	default:
		return errors.New("mimeType " + mimeType + " not supported")
	}
}

func saveImage(w http.ResponseWriter, image multipart.File, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(image)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = ioutil.WriteFile("./storage/uploads/"+handle.Filename, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func imageProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	filename := r.URL.Query()["image"][0]

	p, err := processing.New("storage/uploads/" + filename)
	if err != nil {
		log.Printf("Error loading image for processing: %v", err)
	}
	err = p.SetEffect(processing.EffectLignify)
	if err != nil {
		log.Printf("Error setting effect: %v", err)
	}

	log.Println("Processing...")
	err = p.ProcessImage()
	if err != nil {
		log.Printf("Error processing image: %v", err)
	}

	processedImageFilename, err := p.GetProcessedImageFilename()
	if err != nil {
		log.Printf("Error retrieving processed image filename: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"filename\": \"" + processedImageFilename + "\"}"))

	log.Println("done!")
	fmt.Fprint(w)
}
