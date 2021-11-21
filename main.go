package main

import (
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// Serve up the root html doc
		http.ServeFile(rw, r, "./main.html")
		// io.WriteString(rw, "<!DOCTYPE html><html><h1 style='color: green'>Hello!</h1></html>")
		return
	})

	log.Fatal(http.ListenAndServe(":9090", nil))
}
