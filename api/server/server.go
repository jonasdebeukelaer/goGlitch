package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"text/template"
)

type mainPageVariables struct {
}

type workPageVariables struct {
	Filename string
}

// Serve serves the server
func Serve(port string) {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/upload_image", imageUploadHandler)
	http.HandleFunc("/work", workHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("web/resources/fonts/"))))
	http.Handle("/source_image/", http.StripPrefix("/source_image/", http.FileServer(http.Dir("uploads/"))))

	log.Fatal(http.ListenAndServe(port, nil))
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

	fmt.Print("Uploading image ", handle.Filename, "...")

	http.Redirect(w, r, "/work?image="+handle.Filename, http.StatusSeeOther)
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
