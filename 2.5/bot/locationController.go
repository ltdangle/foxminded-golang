package bot

type WeatherApiInterface interface {
	GetWeather(lat float64, long float64) (string, error)
}

// locationController reacts to geolocation sent by user.
type locationController struct {
	storage StoreInterface
	log     LoggerInterface
}

func NewLocationController(storage StoreInterface, log LoggerInterface) *locationController {
	return &locationController{log: log, storage: storage}
}

func (c *locationController) CreateReply(update Update) *BotMessage {
	sub := c.storage.GetSubscriptionByChatId(update.Message.Chat.ID)

	// if subscription not found, create it
	if sub == nil {
		sub = NewSubscription()
	}

	sub.ChatId = update.Message.Chat.ID
	sub.Latitude = update.Message.Location.Latitude
	sub.Longitude = update.Message.Location.Longitude
	sub.FirsName = update.Message.Chat.FirstName
	sub.Username = update.Message.Chat.Username

	err := c.storage.SaveLocation(sub)
	if err != nil {
		c.log.Warn(err)
		return NewBotMessage("Could not save subscription.")
	}

	return NewBotMessage("We have your geolocation! Now send \"/sub xx:xx:xx\" to subscribe to weather updates.")
}

func (c *locationController) MatchMsg(update Update) bool {
	return update.Message.Location != nil
}
