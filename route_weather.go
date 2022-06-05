package main

import "net/http"

func getWeather(writer http.ResponseWriter, request *http.Request) {
	session, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		user, err := session.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		weather, err := user.GetWeather()
		if err != nil {
			warning(err, "Cannot get weather")
		}
		generateHTML(writer, &weather, "layout", "private.navbar", "weather")
	}
}
