package main

import (
	"net/http"

	"github.com/dkruggel/go-services/data"
)

// GET /err?msg=
// Shows the error page
func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "index")
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "index")
	} else {
		weather, err := data.GetWeather()
		if err != nil {
			error_message(writer, request, "Cannot get weather")
		}
		generateHTML(writer, &weather, "layout", "private.navbar", "home", "weather")
	}
}
