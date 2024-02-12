package bot

import (
	"errors"
	"fmt"
	"time"
)

type CountryData struct {
	Iso   string
	Name  string
	Emoji string
}

// Holidays api interface.
type holidaysInterface interface {
	GetHoliday(date time.Time, country string) ([]string, error)
}

type countryController struct {
	countries []CountryData
	holidays  holidaysInterface
	log       loggerInterface
}

func NewCountryController(holidays holidaysInterface, log loggerInterface) *countryController {
	return &countryController{holidays: holidays, log: log}
}
func (c *countryController) CreateReply(msg Message) *BotMessage {
	botMessage := &BotMessage{}
	emoji := *msg.Text
	country, err := c.countryByEmoji(emoji)
	if err != nil {
		c.log.Warn(err)
	}

	holidays, err := c.holidays.GetHoliday(time.Now(), country.Iso)
	if err != nil {
		botMessage.Text = "Could not get holidays for today for " + country.Iso + ". Sorry...."
		c.log.Warn(err)
	} else {
		botMessage.Text = c.holidayMessage(country.Name, holidays)
	}
	return botMessage
}

func (c *countryController) MatchMsg(msg Message) bool {
	if msg.Text == nil {
		return false
	}

	for _, country := range c.countries {
		if *msg.Text == country.Emoji {
			return true
		}
	}

	return false
}

func (c *countryController) countryByEmoji(emoji string) (CountryData, error) {
	for _, country := range c.countries {
		if emoji == country.Emoji {
			return country, nil
		}
	}
	return CountryData{}, errors.New("can't find country by emoji " + emoji)
}

func (c *countryController) holidayMessage(country string, holidays []string) string {
	var msg string

	if len(holidays) == 0 {
		return fmt.Sprintf("No holidays in %s today.", country)
	}

	msg = fmt.Sprintf("Holidays in %s today:\n", country)
	for _, h := range holidays {
		msg = fmt.Sprintf("%s- %s\n", msg, h)
	}

	return msg
}

func (c *countryController) AddCountry(data CountryData) {
	c.countries = append(c.countries, data)
}
