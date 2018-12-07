package server

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

type pageVariables struct {
	Date string
	Time string
}

// Serve serves the server
func Serve(port string) {
	http.HandleFunc("/", defaultHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("web/resources/fonts/"))))
	log.Fatal(http.ListenAndServe(port, nil))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	testPageVariables := pageVariables{
		Date: now.Format("01/01/2018"),
		Time: now.Format("10:10:10"),
	}

	t, err := template.ParseFiles("web/templates/main.html")
	if err != nil {
		log.Fatalf("template parsing error: %v", err)
	}

	err = t.Execute(w, testPageVariables)
	if err != nil {
		log.Fatalf("populating the template failed: %v", err)
	}
}
