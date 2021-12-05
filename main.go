package main

import (
	"log"
	"net/http"

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
	r.Handle("/notes/{date}", NotesHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":9090", r))
}

var LoadMainPage = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi! My name is David Kruggel and I am a software engineer specializing in .NET and Go."))
})

var StatusCheck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running."))
})

var HomeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	weathertext := weatherservice.GetWeather()
	w.Write([]byte(weathertext))
})

var NotesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notetext := standupnotesservice.GetNote(vars["date"])

	w.Write([]byte(notetext))
})
