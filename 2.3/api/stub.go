package api

import (
	"strconv"
	"time"
)

// stupApi holiday interface implementation.
type stubApi struct{}

func NewStubApi() *stubApi {
	return &stubApi{}
}

func (api *stubApi) Holiday(date time.Time, country string) ([]string, error) {
	holidays := []string{
		strconv.Itoa(date.Day()) + " holiday one " + country,
		strconv.Itoa(date.Day()) + " holiday two" + country,
	}
	return holidays, nil
}
