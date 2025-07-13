package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/simulate", messageHandler)
	http.HandleFunc("/clear", clearHandler)
	http.HandleFunc("/pay", payHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(SimulationResultPage))
}

func clearHandler(w http.ResponseWriter, _ *http.Request) {
	// Send empty string to remove the div
	w.Write([]byte(""))
}

func payHandler(w http.ResponseWriter, _ *http.Request) {
	// Send empty string to remove the div
	w.Write([]byte(NewRowSimulationTable))
}
