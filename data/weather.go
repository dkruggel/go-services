package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Current WeatherNow
	Hourly  []WeatherNow
	Daily   WeatherDaily
}

type WeatherNow struct {
	Time      int64 `json:"dt"`
	Temp      float64
	RealFeel  float64     `json:"feels_like"`
	Condition []Condition `json:"weather"`
}

type WeatherDaily struct {
}

type Condition struct {
	ID   int
	Main string `json:"main"`
}

func (user *User) GetWeather() (weather Weather, err error) {
	client := &http.Client{}
	addr := fmt.Sprintf("%s%s%s", "https://api.openweathermap.org/data/2.5/onecall", "?&lat=38.838159&lon=-90.724872&exclude=minutely&appid=", os.Getenv("WEATHER_API_KEY"))
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		fmt.Printf("FAIL")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("FAIL")
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	fmt.Printf(string(bodyBytes)[0:750])

	json.Unmarshal(bodyBytes, &weather)

	return weather, err
}

func (weather WeatherNow) TempCelcius() (temp string) {
	return fmt.Sprintf("%.2f", weather.Temp-273.15)
}

func (weather WeatherNow) RealFeelCelcius() (temp string) {
	return fmt.Sprintf("%.2f", weather.RealFeel-273.15)
}

func (weather WeatherNow) DisplayTime() (t int64) {
	return int64(time.Unix(weather.Time, 0).Hour())
}
