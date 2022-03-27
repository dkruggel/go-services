package main

import (
	"net/http"
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
		generateHTML(writer, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, "layout", "private.navbar", "index")
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, "layout", "private.navbar", "home")
	}
}
