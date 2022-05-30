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
	session, err := session(writer, request)
	if err != nil {
		generateHTML(writer, request, "layout", "public.navbar", "index")
	} else {
		user, err := session.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		notes, err := user.Notes()
		if err != nil {
			error_message(writer, request, "Cannot retrieve notes")
		} else {
			generateHTML(writer, &notes, "layout", "private.navbar", "notes")
		}
	}
}
