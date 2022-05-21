package main

import (
	"fmt"
	"net/http"
	"strings"
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

func editNote(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		vals := request.URL.Query()
		uuid := vals.Get("id")
		note, err := data.NoteByUUID(uuid)
		if err != nil {
			error_message(writer, request, "Cannot find note")
		} else {
			generateHTML(writer, &note, "layout", "private.navbar", "edit.note")
		}
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
		vals := request.URL.Query()
		uuid := vals.Get("id")
		if _, err := user.CreateNote(date, yesterday, today, gobacks, impediments); err != nil {
			danger(err, "Cannot create note")
		}
		redir := fmt.Sprintf("/note?id=%s", uuid)
		http.Redirect(writer, request, redir, http.StatusFound)
	}
}

func getNotes(writer http.ResponseWriter, request *http.Request) {
	session, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}
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

func updateNote(writer http.ResponseWriter, request *http.Request) {
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
		query := strings.Split(request.Referer(), "?")[1]
		uuid := query[3:]
		date := request.Form["date"][0]
		yesterday := request.Form["yesterday"][0]
		today := request.Form["today"][0]
		gobacks := request.Form["gobacks"][0]
		impediments := request.Form["impediments"][0]
		if _, err := user.UpdateNote(uuid, date, yesterday, today, gobacks, impediments); err != nil {
			danger(err, "Cannot create note")
		}
		redir := fmt.Sprintf("/note?id=%s", uuid)
		http.Redirect(writer, request, redir, http.StatusFound)
	}
}

func deleteNote(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("About to create a session\n")
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		vals := request.URL.Query()
		uuid := vals.Get("id")
		fmt.Printf("%s", uuid)
		err = user.DeleteNote(uuid)
		if err != nil {
			danger(err, "Cannot delete note")
		}
		http.Redirect(writer, request, "/notes", http.StatusFound)
	}
}
