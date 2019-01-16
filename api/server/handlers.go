package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/jonasdebeukelaer/goGlitch/processing"
)

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
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	processedImageFilename, err := prepareAndProcessImage(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
	} else {
		w.Write([]byte("{\"filename\": \"" + processedImageFilename + "\"}"))
	}

	fmt.Fprint(w)
}

func prepareAndProcessImage(r *http.Request) (string, error) {
	filename := r.URL.Query()["image"][0]
	effectLayerList, err := parseEffectLayerList(r.Body)
	if err != nil {
		log.Printf("error getting reading effect layers: \n\t%v", err.Error())
	}

	p, err := runProcessing(filename, effectLayerList)
	if err != nil {
		log.Printf("error during image processing: \n\t%v", err.Error())
	}

	return p.GetProcessedImageFilename()
}

func parseEffectLayerList(bodyIo io.ReadCloser) ([]*processing.EffectLayer, error) {
	body, err := ioutil.ReadAll(bodyIo)
	if err != nil {
		return nil, errors.New("error reading body content")
	}

	var effectLayerList []*processing.EffectLayer
	err = json.Unmarshal(body, &effectLayerList)
	if err != nil {
		return nil, errors.New("error parsing body content")
	}

	return effectLayerList, nil
}

func runProcessing(filename string, layers []*processing.EffectLayer) (*processing.Processor, error) {
	log.Println("Processing...")
	p, err := processing.New("storage/uploads/" + filename)
	if err != nil {
		return nil, fmt.Errorf("error loading image for processing: %v", err)
	}

	err = p.ProcessImage(layers)
	if err != nil {
		return nil, fmt.Errorf("error processing image: %v", err)
	}
	log.Println("done!")

	return p, nil
}

func effectOptionHandler(w http.ResponseWriter, r *http.Request) {

	listBytes, err := json.Marshal(processing.EffectList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to parse effects list")
		fmt.Fprintf(w, "failed to parse effects list")
	}

	_, err = w.Write(listBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to write effects list")
		fmt.Fprintf(w, "failed to write effects list")
	}
}
