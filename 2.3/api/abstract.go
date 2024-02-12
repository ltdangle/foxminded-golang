package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AbstractApi holiday result.
type AbstractApiHoliday struct {
	Name        string `json:"name"`
	Name_local  string `json:"name_local"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Country     string `json:"country"`
	Location    string `json:"location"`
	Type        string `json:"type"`
	Date        string `json:"date"`
	Date_year   string `json:"date_year"`
	Date_month  string `json:"date_month"`
	Date_day    string `json:"date_day"`
	Week_day    string `json:"week_day"`
}

// AbstractApi holiday interface implementation.
type abstractApi struct {
	url    string
	apiKey string
}

func NewAbstractApi(url string, apiKey string) *abstractApi {
	return &abstractApi{url: url, apiKey: apiKey}
}

func (api *abstractApi) GetHoliday(date time.Time, country string) ([]string, error) {
	year := date.Year()
	month := date.Month()
	day := date.Day()

	url := fmt.Sprintf("%s?api_key=%s&country=%s&year=%d&month=%d&day=%d", api.url, api.apiKey, country, year, month, day)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse []AbstractApiHoliday
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	var holidays []string
	for _, holiday := range restResponse {
		holidays = append(holidays, holiday.Name)
	}

	return holidays, nil
}
