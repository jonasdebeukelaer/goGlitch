package server

import (
	"fmt"
	"log"
	"net/http"
)

// Serve serves the server
func Serve(port string) {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/upload_image", imageUploadHandler)
	http.HandleFunc("/work", workHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("web/resources/fonts/"))))
	http.Handle("/source_image/", http.StripPrefix("/source_image/", http.FileServer(http.Dir("uploads/"))))

	fmt.Println("Running server on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
