package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/joho/godotenv"
)

type Weather struct {
	Current WeatherNow     `json:"current"`
	Daily   []WeatherDaily `json:"daily"`
}

type WeatherRain struct {
	LastHour float64 `json:"1h"`
}

type WeatherSnow struct {
	LastHour float64 `json:"1h"`
}

type WeatherInfo struct {
	Id   float64 `json:"id"`
	Main float64 `json:"main"`
	Desc float64 `json:"description"`
	//Icon float64 `json:"icon"`
}

type WeatherTemp struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type WeatherDaily struct {
	Date    float64     `json:"dt"`
	Sunrise float64     `json:"sunrise"`
	Sunset  float64     `json:"sunset"`
	Temp    WeatherTemp `json:"temp"`
}

type WeatherNow struct {
	Sunrise     float64     `json:"sunrise"`
	Sunset      float64     `json:"sunset"`
	CurrentTemp float64     `json:"temp"`
	FeelsLike   float64     `json:"feels_like"`
	Rain        WeatherRain `json:"rain"`
	Snow        WeatherSnow `json:"snow"`
}

func (*Weather) GetCelcius(incoming float64, decPlaces int) float64 {
	x := math.Pow(10, float64(decPlaces))
	final := math.Round((incoming-273.15)*x) / x
	return final
}

func GetLocalTime(inc float64) string {
	loc, _ := time.LoadLocation("America/Chicago")
	return time.Unix(int64(inc), 0).UTC().In(loc).Format("3:04pm")
}

func (weather *Weather) DisplayCurrentTemp() string {
	return fmt.Sprintf("%.2f\u2103", weather.Current.CurrentTemp-273.15)
}

func (weather *Weather) DisplayFeelsLike() string {
	return fmt.Sprintf("%.2f\u2103", weather.Current.FeelsLike-273.15)
}

func GetWeather() (w Weather, err error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

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
		if err := dec.Decode(&w); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		for i := range w.Daily {
			w.Daily[i].Temp.Day = w.GetCelcius(w.Daily[i].Temp.Day, 2)
			w.Daily[i].Temp.Min = w.GetCelcius(w.Daily[i].Temp.Min, 2)
			w.Daily[i].Temp.Max = w.GetCelcius(w.Daily[i].Temp.Max, 2)
			w.Daily[i].Temp.Night = w.GetCelcius(w.Daily[i].Temp.Night, 2)
			w.Daily[i].Temp.Eve = w.GetCelcius(w.Daily[i].Temp.Eve, 2)
			w.Daily[i].Temp.Morn = w.GetCelcius(w.Daily[i].Temp.Morn, 2)
		}

		// text = fmt.Sprintf("Current: %.2f\u2103\nReal Feel: %.2f\u2103\nSunrise: %s\nSunset: %s\n", GetCelcius(c.Current.CurrentTemp),
		// 	GetCelcius(c.Current.FeelsLike), GetLocalTime(c.Current.Sunrise),
		// 	GetLocalTime(c.Current.Sunset))
	}

	return w, err
}
