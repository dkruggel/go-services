package goservices

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	lat := fmt.Sprintf("%.6f", location.Lat)
	lon := fmt.Sprintf("%.6f", location.Lng)

	// Get weather
	// Retrieve environment variable that holds the openweather API key
	apikey := os.Getenv("WEATHER_API_KEY")

	// Build the api call to openweather
	uri := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&appid=%s", lat, lon, apikey)
	resp, err := http.Get(uri)

	if err != nil {
		fmt.Printf("%s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("%s", err)
	}

	dec := json.NewDecoder(strings.NewReader(string(body)))

	for {
		var c WeatherNow
		if err := dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%.2f\u2103\n%.2f\u2103\n%s\n%s\n", GetCelcius(c.Current.CurrentTemp),
			GetCelcius(c.Current.FeelsLike), GetLocalTime(c.Current.Sunrise),
			GetLocalTime(c.Current.Sunset))
	}

	w.Write([]byte(body))
})

func GetCelcius(incoming float64) float64 {
	return incoming - 273.15
}

func GetLocalTime(inc float64) string {
	loc, _ := time.LoadLocation("America/Chicago")
	return time.Unix(int64(inc), 0).UTC().In(loc).Format("3:04pm")
}
