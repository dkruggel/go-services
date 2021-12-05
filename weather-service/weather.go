package goservices

type WeatherNow struct {
	Current WeatherInfo `json:"current"`
}

type WeatherRain struct {
	LastHour float64 `json:"1h"`
}

type WeatherSnow struct {
	LastHour float64 `json:"1h"`
}

type WeatherInfo struct {
	Sunrise     float64     `json:"sunrise"`
	Sunset      float64     `json:"sunset"`
	CurrentTemp float64     `json:"temp"`
	FeelsLike   float64     `json:"feels_like"`
	Rain        WeatherRain `json:"rain"`
	Snow        WeatherSnow `json:"snow"`
}
