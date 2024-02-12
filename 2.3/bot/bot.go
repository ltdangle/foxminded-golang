package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger interface.
type loggerInterface interface {
	Info(args ...interface{})
	Warn(args ...interface{})
}

// Holidays api interface.
type holidaysInterface interface {
	GetHoliday(date time.Time, country string) ([]string, error)
}

type CountryData struct {
	Iso   string
	Name  string
	Emoji string
}

// Bot application.
type bot struct {
	botUrl    string
	countries []CountryData
	holidays  holidaysInterface
	log       loggerInterface
}

func NewBot(botUrl string, holidays holidaysInterface, log loggerInterface) *bot {
	return &bot{botUrl: botUrl, holidays: holidays, log: log}
}

// Main bot loop.
func (bot *bot) Run() {
	offset := 0
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		updates, err := bot.getUpdates(offset)
		if err != nil {
			log.Warn("Smth went wrong: ", err.Error())
		}

		for _, update := range updates {
			err = bot.respond(update)
			if err != nil {
				log.Warn("Smth went wrong: ", err.Error())
			}

			offset = update.UpdateId + 1
		}

		log.Info(updates)
	}
}

// Query bot api for updates.
func (bot *bot) getUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(bot.botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

// Respond to incoming messages (updates).
func (bot *bot) respond(update Update) error {

	respondMessage := bot.runCommand(update.Message.Text)
	respondMessage.ChatId = update.Message.Chat.ChatId

	buf, err := json.Marshal(respondMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(bot.botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}

// Get the type of incoming message.
func (bot *bot) msgType(msgIn string) string {
	for _, country := range bot.countries {
		if msgIn == country.Emoji {
			return "{country}"
		}
	}
	return msgIn
}

func (bot *bot) runCommand(msgIn string) *BotMessage {
	botMessage := &BotMessage{}

	switch bot.msgType(msgIn) {
	case "/start":
		var kb [][]string
		botMessage.Text = "Pick a country..."
		kb = append(kb, bot.countryNames())
		botMessage.ReplyMarkup = &ReplyKeyboardMarkup{
			Keyboard:        kb,
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}
	case "{country}":
		emoji := msgIn
		country, err := bot.countryByEmoji(emoji)
		if err != nil {
			bot.log.Warn(err)
		}

		holidays, err := bot.holidays.GetHoliday(time.Now(), country.Iso)
		if err != nil {
			botMessage.Text = "Could not get holidays for today for " + country.Iso + ". Sorry...."
			bot.log.Warn(err)
		} else {
			botMessage.Text = bot.holidayMessage(country.Name, holidays)
		}
	default:
		botMessage.Text = "Command not found..."
	}

	return botMessage
}

func (bot *bot) holidayMessage(country string, holidays []string) string {
	var msg string

	if len(holidays) == 0 {
		return fmt.Sprintf("No holidyas in %s today.", country)
	}

	msg = fmt.Sprintf("Holidays in %s today:\n", country)
	for _, h := range holidays {
		msg = fmt.Sprintf("%s- %s\n", msg, h)
	}

	return msg
}

func (bot *bot) countryByEmoji(emoji string) (CountryData, error) {
	for _, country := range bot.countries {
		if emoji == country.Emoji {
			return country, nil
		}
	}
	return CountryData{}, errors.New("can't find country by emoji " + emoji)
}

func (bot *bot) countryNames() []string {
	var names []string
	for _, country := range bot.countries {
		names = append(names, country.Emoji)
	}
	return names
}

func (bot *bot) AddCountry(data CountryData) {
	bot.countries = append(bot.countries, data)
}
