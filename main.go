package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	standupnotesservice "github.com/dkruggel/go-services/standup-notes-service"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Load main page
	r.HandleFunc("/", index)

	// Load error page
	r.HandleFunc("/err", err)

	// Auth pages
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/signup_account", signupAccount)
	r.HandleFunc("/authenticate", authenticate)

	// Weather
	r.HandleFunc("/weather", getWeather) // Get weather data

	// Stand up notes
	r.HandleFunc("/notes", getNotes)                          // Get all of the notes
	r.HandleFunc("/note", getNote)                            // Get one specific note
	r.HandleFunc("/note/new", newNote)                        // Get new note page
	r.HandleFunc("/note/create", createNote)                  // Create new note
	r.HandleFunc("/note/edit", editNote)                      // Edit existing note
	r.HandleFunc("/note/update", updateNote)                  // Update existing note
	r.HandleFunc("note/delete", deleteNote).Methods("DELETE") // Delete existing note - TODO: FIX
	r.Handle("/note/{date}", NotesHandler).Methods("GET")
	r.Handle("/note/{date}", NotesHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}

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
	case http.MethodGet:
		note, err := standupnotesservice.GetNote(vars["date"])

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	case http.MethodPost:
		w.Write([]byte("Saved"))
	case http.MethodPut:
	case http.MethodDelete:
		w.Write([]byte("Not Implemented Yet"))
	}
})
