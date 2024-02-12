package bot

import (
	"encoding/json"
	"time"
)

type weatherPusher struct {
	store      StoreInterface
	weatherApi WeatherApiInterface
	log        LoggerInterface
}

type botInterface interface {
	respond(msgJson []byte) error
}

func NewWeatherPusher(store StoreInterface, weather WeatherApiInterface, log LoggerInterface) *weatherPusher {
	return &weatherPusher{store: store, weatherApi: weather, log: log}
}

func (p *weatherPusher) Push(bot botInterface, t time.Time) {
	// Get all subscriptions matching current time.
	timeStr := t.UTC().Format("15:04:05")

	subs, err := p.store.FindSubsBySchedule(timeStr)
	if err != nil {
		p.log.Warn("Smth went wrong: ", err.Error())
	}

	p.log.Info("Subscriptions at " + timeStr)
	p.log.Info(subs)

	// For each matching subsription, push current weather
	for _, sub := range subs {
		var botMessage *BotMessage

		weather, err := p.weatherApi.GetWeather(sub.Latitude, sub.Longitude)
		if err != nil {
			p.log.Warn(err)
			botMessage = NewBotMessage("Could not retrieve weather...sorry!")
		} else {
			botMessage = NewBotMessage(weather)
		}

		botMessage.ChatId = sub.ChatId

		msgJson, err := json.Marshal(botMessage)
		if err != nil {
			p.log.Warn("Smth went wrong: ", err.Error())
			botMessage.Text = "Could not parse weather"
		}

		err = bot.respond(msgJson)
		if err != nil {
			p.log.Warn("Smth went wrong: ", err.Error())
		}
	}
}
