package server

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/jonasdebeukelaer/goGlitch/processing"
	"github.com/jonasdebeukelaer/goGlitch/processing/processor"
)

type mainPageVariables struct {
}

type workPageVariables struct {
	Filename string
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/main.html")
	if err != nil {
		log.Fatalf("template parsing error: %v", err)
	}

	testPageVariables := mainPageVariables{}

	err = t.Execute(w, testPageVariables)
	if err != nil {
		log.Fatalf("populating the template failed: %v", err)
	}
}

func imageUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	image, handle, err := r.FormFile("picture")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer image.Close()
	fmt.Print("Uploading image ", handle.Filename, "...")

	err = checkFileType(handle)
	if err != nil {
		fmt.Println("error!\nImage could not be uploaded: %", err)
		return
	}

	saveImage(w, image, handle)
	fmt.Println("success!")

	http.Redirect(w, r, "/work?image="+handle.Filename, http.StatusSeeOther)
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

	err = ioutil.WriteFile("./uploads/"+handle.Filename, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/work.html")
	if err != nil {
		log.Fatalf("work template parsing error: %v", err)
	}

	r.ParseForm()
	workPageVariables := workPageVariables{
		Filename: r.Form["image"][0],
	}

	err = t.Execute(w, workPageVariables)
	if err != nil {
		log.Fatalf("populating the template failed: %v", err)
	}
}

func imageProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	filename := r.URL.Query()["image"][0]

	p, err := processor.New("uploads/" + filename)
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
	w.Write([]byte("{\"imageURL\": \"" + processedImageFilename + "\"}"))

	log.Println("done!")
	fmt.Fprint(w)
}
