package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mockup-billing-engine/repo"
	"github.com/mockup-billing-engine/usecase"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {

	//init DB
	dbClient := repo.Init()
	defer dbClient.CloseDB()
	//init Usecase
	uc := usecase.Init(dbClient)

	// fs := http.FileServer(http.Dir("static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/simulate", uc.SimulateHandler)
	http.HandleFunc("/clear", uc.ClearHandler)
	http.HandleFunc("/pay", uc.PayHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}
