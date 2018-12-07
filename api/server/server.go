package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"text/template"
	"time"
)

type pageVariables struct {
	Date string
}

// Serve serves the server
func Serve(port string) {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/upload_image", imageUploadHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("web/resources/fonts/"))))
	log.Fatal(http.ListenAndServe(port, nil))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/main.html")
	if err != nil {
		log.Fatalf("template parsing error: %v", err)
	}

	now := time.Now()
	testPageVariables := pageVariables{
		Date: now.Format("01/01/2018 10:01:10"),
	}

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
	fmt.Print("Uploading image ", r.Form["caption"], "...")

	image, handle, err := r.FormFile("pic")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer image.Close()

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		saveImage(w, image, handle)
	case "image/png":
		saveImage(w, image, handle)
	default:
		fmt.Println("file type ", mimeType, " not supported")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	fmt.Println("success!")
}
