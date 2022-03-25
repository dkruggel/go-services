package weatherservice

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

func (*Weather) GetCelcius(incoming float64) float64 {
	return incoming - 273.15
}

func GetLocalTime(inc float64) string {
	loc, _ := time.LoadLocation("America/Chicago")
	return time.Unix(int64(inc), 0).UTC().In(loc).Format("3:04pm")
}

func GetWeather() Weather {
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

	var c Weather

	for {
		if err := dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		c.Current.CurrentTemp = c.GetCelcius(c.Current.CurrentTemp)
		c.Current.FeelsLike = c.GetCelcius(c.Current.FeelsLike)

		for i := range c.Daily {
			c.Daily[i].Temp.Day = c.GetCelcius(c.Daily[i].Temp.Day)
			c.Daily[i].Temp.Min = c.GetCelcius(c.Daily[i].Temp.Min)
			c.Daily[i].Temp.Max = c.GetCelcius(c.Daily[i].Temp.Max)
			c.Daily[i].Temp.Night = c.GetCelcius(c.Daily[i].Temp.Night)
			c.Daily[i].Temp.Eve = c.GetCelcius(c.Daily[i].Temp.Eve)
			c.Daily[i].Temp.Morn = c.GetCelcius(c.Daily[i].Temp.Morn)
		}

		// text = fmt.Sprintf("Current: %.2f\u2103\nReal Feel: %.2f\u2103\nSunrise: %s\nSunset: %s\n", GetCelcius(c.Current.CurrentTemp),
		// 	GetCelcius(c.Current.FeelsLike), GetLocalTime(c.Current.Sunrise),
		// 	GetLocalTime(c.Current.Sunset))
	}

	return c
}
