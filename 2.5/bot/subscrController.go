package bot

import (
	"strings"
	"time"
)

// subscrController creates new weather push subscription at specified time
type subscrController struct {
	// Subscription time.
	schedule string
	store    StoreInterface
	log      LoggerInterface
}

func NewSubscrController(store StoreInterface, log LoggerInterface) *subscrController {
	return &subscrController{store: store, log: log}
}

func (c *subscrController) CreateReply(update Update) *BotMessage {
	sub := c.store.GetSubscriptionByChatId(update.Message.Chat.ID)
	if sub == nil {
		return NewBotMessage("Subscription request not found. Did you already send your geolocation?")
	}

	if sub.Latitude == 0 && sub.Longitude == 0 {
		return NewBotMessage("Subscription request not found. Did you already send your geolocation?")
	}

	err := c.store.SaveSchedule(sub.Id, c.schedule)
	if err != nil {
		c.log.Warn(err)
		return NewBotMessage("Could not save subscription schedule.")
	}

	return NewBotMessage("You will receive weather updates each day @ " + c.schedule + " UTC")
}

func (c *subscrController) MatchMsg(update Update) bool {
	if update.Message.Text == nil {
		return false
	}

	str := *update.Message.Text
	if !strings.HasPrefix(str, "/sub ") {
		return false
	}

	// Validate time.
	splitStr := strings.Split(str, " ")
	subTime := splitStr[1]
	_, err := time.Parse("15:04:05", subTime)
	if err != nil {
		return false
	}

	c.schedule = subTime
	return true
}
