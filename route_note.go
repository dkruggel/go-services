package main

import (
	"net/http"
	"time"

	"github.com/dkruggel/go-services/data"
)

func newNote(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		var date string
		_, err := data.NoteByDate(formatDate(time.Now()))
		_, exists := data.NoteByDate(formatDate(time.Now().AddDate(0, 0, 1)))
		if err != nil {
			date = formatDate(time.Now())
		} else if time.Now().Weekday() == time.Friday {
			date = formatDate(time.Now().AddDate(0, 0, 3))
		} else if exists == nil {
			date = ""
		} else {
			date = formatDate(time.Now().AddDate(0, 0, 1))
		}
		generateHTML(writer, date, "layout", "private.navbar", "new.note")
	}
}

func getNote(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	note, err := data.NoteByUUID(uuid)
	if err != nil {
		error_message(writer, request, "Cannot find note")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, &note, "layout", "public.navbar", "public.note")
		} else {
			generateHTML(writer, &note, "layout", "private.navbar", "private.note")
		}
	}
}

func createNote(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		date := request.PostFormValue("date")
		yesterday := request.PostFormValue("yesterday")
		today := request.PostFormValue("today")
		gobacks := request.PostFormValue("gobacks")
		impediments := request.PostFormValue("impediments")
		if _, err := user.CreateNote(date, yesterday, today, gobacks, impediments); err != nil {
			danger(err, "Cannot create note")
		}
		http.Redirect(writer, request, "/notes", http.StatusFound)
	}
}

func getNotes(writer http.ResponseWriter, request *http.Request) {
	notes, err := data.Notes()
	if err != nil {
		error_message(writer, request, "Cannot retrieve notes")
	} else {
		_, err := session(writer, request)
		if err != nil {
			http.Redirect(writer, request, "/", http.StatusFound)
		} else {
			generateHTML(writer, &notes, "layout", "private.navbar", "notes")
		}
	}
}
