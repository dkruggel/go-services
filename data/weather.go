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
	Daily   []WeatherDaily
}

type WeatherNow struct {
	Time      int64 `json:"dt"`
	Temp      float64
	RealFeel  float64     `json:"feels_like"`
	Condition []Condition `json:"weather"`
}

type WeatherDaily struct {
	Date      int64 `json:"dt"`
	Sunrise   int64
	Sunset    int64
	Temp      Temperatures
	Condition []Condition `json:"weather"`
}

type Condition struct {
	ID   int
	Main string `json:"main"`
	Icon string
}

type Temperatures struct {
	Day     float64
	Min     float64
	Max     float64
	Night   float64
	Evening float64 `json:"eve"`
	Morning float64 `json:"morn"`
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

	fmt.Printf(string(bodyBytes))

	json.Unmarshal(bodyBytes, &weather)

	return weather, err
}

func (weather WeatherNow) TempCelcius() (temp string) {
	return fmt.Sprintf("%.2f", weather.Temp-273.15)
}

func (temperature Temperatures) MaxTempCelcius() (temp string) {
	return fmt.Sprintf("%.2f", temperature.Max-273.15)
}

func (temperature Temperatures) MinTempCelcius() (temp string) {
	return fmt.Sprintf("%.2f", temperature.Min-273.15)
}

func (weather WeatherNow) RealFeelCelcius() (temp string) {
	return fmt.Sprintf("%.2f", weather.RealFeel-273.15)
}

func (weather WeatherNow) DisplayTime() (t string) {
	dt := int64(time.Unix(weather.Time, 0).Local().Hour())
	if dt > 12 {
		return fmt.Sprintf("%d:00p", dt%12)
	} else if dt == 0 {
		return "12:00a"
	} else {
		return fmt.Sprintf("%d:00a", dt)
	}
}

func (weather WeatherDaily) DisplayDate() (t string) {
	_, m, d := time.Unix(weather.Date, 0).Date()
	w := time.Unix(weather.Date, 0).Weekday()
	return fmt.Sprintf("%v, %v %d", w, m, d)
}

func (condition Condition) DisplayIcon() (s string) {
	return fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", condition.Icon)
}
