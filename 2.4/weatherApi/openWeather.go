package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type weatherApi struct {
	url    string
	apiKey string // https://api.openweathermap.org/data/2.5/weather
}

type openWeatherResponse struct {
	Base   string `json:"base,omitempty"`
	Clouds struct {
		All float64 `json:"all,omitempty"`
	} `json:"clouds,omitempty"`
	Cod   float64 `json:"cod,omitempty"`
	Coord struct {
		Lat float64 `json:"lat,omitempty"`
		Lon float64 `json:"lon,omitempty"`
	} `json:"coord,omitempty"`
	Dt   float64 `json:"dt,omitempty"`
	ID   float64 `json:"id,omitempty"`
	Main struct {
		FeelsLike float64 `json:"feels_like,omitempty"`
		Humidity  float64 `json:"humidity,omitempty"`
		Pressure  float64 `json:"pressure,omitempty"`
		Temp      float64 `json:"temp,omitempty"`
		TempMax   float64 `json:"temp_max,omitempty"`
		TempMin   float64 `json:"temp_min,omitempty"`
	} `json:"main,omitempty"`
	Name string `json:"name,omitempty"`
	Sys  struct {
		Country string  `json:"country,omitempty"`
		ID      float64 `json:"id,omitempty"`
		Sunrise float64 `json:"sunrise,omitempty"`
		Sunset  float64 `json:"sunset,omitempty"`
		Type    float64 `json:"type,omitempty"`
	} `json:"sys,omitempty"`
	Timezone   float64 `json:"timezone,omitempty"`
	Visibility float64 `json:"visibility,omitempty"`
	Weather    []struct {
		Description string  `json:"description,omitempty"`
		Icon        string  `json:"icon,omitempty"`
		ID          float64 `json:"id,omitempty"`
		Main        string  `json:"main,omitempty"`
	} `json:"weather,omitempty"`
	Wind struct {
		Deg   float64 `json:"deg,omitempty"`
		Gust  float64 `json:"gust,omitempty"`
		Speed float64 `json:"speed,omitempty"`
	} `json:"wind,omitempty"`
}

func NewWeatherApi(url string, apiKey string) *weatherApi {
	return &weatherApi{url: url, apiKey: apiKey}
}

func (api *weatherApi) GetWeather(lat float64, long float64) (string, error) {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", api.url, lat, long, api.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r openWeatherResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}

	weather := fmt.Sprintf(`
Weather for %s, %s
Temp: %.0f C
Wind: %.0f km/h
Humidity: %.0f %%
`,
		r.Name, r.Sys.Country, r.Main.Temp, r.Wind.Speed, r.Main.Humidity)

	return weather, err
}
