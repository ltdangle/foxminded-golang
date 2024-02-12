package bot

type WeatherApiInterface interface {
	GetWeather(lat float64, long float64) (string, error)
}

type weatherController struct {
	weatherApi WeatherApiInterface
	log        loggerInterface
}

func NewWeatherController(weatherApi WeatherApiInterface, log loggerInterface) *weatherController {
	return &weatherController{weatherApi: weatherApi, log: log}
}

func (c *weatherController) CreateReply(msg Message) *BotMessage {
	weather, err := c.weatherApi.GetWeather(msg.Location.Latitude, msg.Location.Longitude)
	if err != nil {
		c.log.Warn(err)
		return &BotMessage{Text: "Could not retrive weather, sorry!"}
	}

	return &BotMessage{Text: weather}
}

func (c *weatherController) MatchMsg(msg Message) bool {
	return msg.Location != nil
}
