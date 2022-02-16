package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	standupnotesservice "github.com/dkruggel/go-services/standup-notes-service"
	weatherservice "github.com/dkruggel/go-services/weather-service"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Load main page
	r.Handle("/", LoadMainPage).Methods("GET")

	r.Handle("/status", StatusCheck).Methods("GET")

	// Weather
	r.Handle("/weather", HomeHandler).Methods("GET")

	// Stand up notes
	r.Handle("/notes", NoteHandler).Methods("GET")
	r.Handle("/note/{date}", NotesHandler).Methods("GET")
	r.Handle("/note/{date}", NotesHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}

var LoadMainPage = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi! My name is David Kruggel and I am a software engineer specializing in .NET framework/core, Go, and React."))
})

var StatusCheck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode("API is up and running.")
})

var HomeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	weathertext := weatherservice.GetWeather()
	json.NewEncoder(w).Encode(weathertext)
})

var NoteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	notelist := standupnotesservice.GetAllNotes()
	json.NewEncoder(w).Encode(notelist)
})

var NotesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	switch r.Method {
	case "GET":
		note := standupnotesservice.GetNote(vars["date"])

		json.NewEncoder(w).Encode(note)
	case "POST":
		w.Write([]byte("Saved"))
	}
})
