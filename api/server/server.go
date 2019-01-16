package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

// Serve serves the server
func Serve(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", defaultHandler)
	mux.HandleFunc("/upload_image", imageUploadHandler)
	mux.HandleFunc("/process_image", imageProcessHandler)
	mux.HandleFunc("/effect_options", effectOptionHandler)

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css/"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("web/js/"))))
	mux.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("node_modules/"))))
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("web/resources/fonts/"))))
	mux.Handle("/source_image/", http.StripPrefix("/source_image/", http.FileServer(http.Dir("storage/uploads/"))))
	mux.Handle("/processed_image/", http.StripPrefix("/processed_image/", http.FileServer(http.Dir("storage/processed_images/"))))

	handler := cors.Default().Handler(mux)

	fmt.Println("Server running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, handler))
}
