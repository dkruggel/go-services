package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/status", StatusCheck).Methods("GET")

	// Home page
	r.Handle("/home", HomeHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":9090", r))
}

var StatusCheck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running."))
})

var HomeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Get location in lat/lon
	location, _ := openstreetmap.Geocoder().Geocode("996 Crestwood Lane, O'Fallon, MO, 63366")
	fmt.Printf("%.6f : %.6f", location.Lat, location.Lng)
	s := fmt.Sprintf("%.6f : %.6f", location.Lat, location.Lng)
	w.Write([]byte(s))

	// Get weather
	
})
