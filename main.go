package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", defaultHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test page")
}
